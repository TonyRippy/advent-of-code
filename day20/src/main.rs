use std::collections::{HashMap, VecDeque};
use std::fmt::{self, Debug, Formatter};
use std::fs::File;
use std::io::BufRead;

#[derive(Clone, Copy, PartialEq, Eq)]
enum Pulse {
    Low,
    High,
}

impl Debug for Pulse {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match self {
            Pulse::Low => write!(f, "-low->"),
            Pulse::High => write!(f, "-high->"),
        }
    }
}

struct Signal {
    pub from: String,
    pub pulse: Pulse,
    pub to: String,
}

impl Debug for Signal {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        write!(f, "{} {:?} {}", self.from, self.pulse, self.to)
    }
}

trait Receiver {
    fn name(&self) -> &str;
    fn add_input(&mut self, name: &str);
    fn add_output(&mut self, name: &str);
    fn receive(&mut self, signal: Signal, queue: &mut VecDeque<Signal>);
}

#[derive(Default)]
struct Output {}

impl Receiver for Output {
    fn name(&self) -> &str {
        "output"
    }

    fn add_input(&mut self, _: &str) {}

    fn add_output(&mut self, _: &str) {
        panic!("Output node cannot have outputs");
    }

    fn receive(&mut self, _: Signal, _: &mut VecDeque<Signal>) {
        // NOOP
    }
}

#[derive(Default)]
struct Broadcast {
    outputs: Vec<String>,
}

impl Receiver for Broadcast {
    fn name(&self) -> &str {
        "broadcaster"
    }

    fn add_input(&mut self, _: &str) {
        panic!("Broadcast node cannot have inputs");
    }

    fn add_output(&mut self, name: &str) {
        self.outputs.push(name.to_string());
    }

    fn receive(&mut self, signal: Signal, queue: &mut VecDeque<Signal>) {
        for output in &self.outputs {
            queue.push_back(Signal {
                from: "broadcaster".to_string(),
                pulse: signal.pulse,
                to: output.clone(),
            });
        }
    }
}

struct FlipFlop {
    name: String,
    on: bool,
    outputs: Vec<String>,
}

impl FlipFlop {
    fn new(name: String) -> Self {
        Self {
            name,
            on: false,
            outputs: Vec::new(),
        }
    }
}

impl Receiver for FlipFlop {
    fn name(&self) -> &str {
        &self.name
    }
    fn add_input(&mut self, _: &str) {}

    fn add_output(&mut self, name: &str) {
        self.outputs.push(name.to_string());
    }

    fn receive(&mut self, signal: Signal, queue: &mut VecDeque<Signal>) {
        assert_eq!(signal.to, self.name);

        if signal.pulse == Pulse::High {
            // ignore and move on
            return;
        }
        let send = match self.on {
            true => {
                self.on = false;
                Pulse::Low
            }
            false => {
                self.on = true;
                Pulse::High
            }
        };
        for output in &self.outputs {
            queue.push_back(Signal {
                from: self.name.clone(),
                pulse: send,
                to: output.clone(),
            });
        }
    }
}

struct Conjunction {
    name: String,
    inputs: HashMap<String, Pulse>,
    outputs: Vec<String>,
}

impl Conjunction {
    fn new(name: String) -> Self {
        Self {
            name,
            inputs: HashMap::new(),
            outputs: Vec::new(),
        }
    }
}

impl Receiver for Conjunction {
    fn name(&self) -> &str {
        &self.name
    }
    fn add_input(&mut self, name: &str) {
        self.inputs.insert(name.to_string(), Pulse::Low);
    }

    fn add_output(&mut self, name: &str) {
        self.outputs.push(name.to_string());
    }

    fn receive(&mut self, signal: Signal, queue: &mut VecDeque<Signal>) {
        assert_eq!(signal.to, self.name);

        // Update memory of the input
        self.inputs.insert(signal.from.clone(), signal.pulse);
        // Check if all inputs are high
        let all_high = self.inputs.values().all(|&p| p == Pulse::High);
        for output in &self.outputs {
            queue.push_back(Signal {
                from: self.name.clone(),
                pulse: if all_high { Pulse::Low } else { Pulse::High },
                to: output.clone(),
            });
        }
    }
}

struct Network {
    nodes: HashMap<String, Box<dyn Receiver>>,
    queue: VecDeque<Signal>,
    presses: usize,
    low_signals: usize,
    high_signals: usize,
}

impl Network {
    fn parse(reader: impl BufRead) -> Self {
        let mut nodes: HashMap<String, Box<dyn Receiver>> = HashMap::new();
        let mut pairs: Vec<(String, String)> = Vec::new();
        for line in reader.lines() {
            let line = line.unwrap();
            let mut parts = line.split(" -> ");
            let input = parts.next().unwrap();
            let node: Box<dyn Receiver> = if input == "broadcaster" {
                Box::<Broadcast>::default()
            } else {
                let mut name = input.to_string();
                name.remove(0); // remove the leading % or &
                match input.chars().next().unwrap() {
                    '%' => Box::new(FlipFlop::new(name)),
                    '&' => Box::new(Conjunction::new(name)),
                    _ => panic!("Unknown node type"),
                }
            };
            let name = node.name().to_string();
            nodes.insert(name.clone(), node);
            for output in parts
                .next()
                .unwrap()
                .split(',')
                .map(|s| s.trim().to_string())
            {
                pairs.push((name.clone(), output));
            }
        }
        for (input, output) in pairs {
            nodes
                .get_mut(&input)
                .expect("unable to find input node")
                .add_output(&output);
            if let Some(node) = nodes.get_mut(&output) {
                node.add_input(&input);
            }
        }
        Self {
            nodes,
            queue: VecDeque::new(),
            presses: 0,
            low_signals: 0,
            high_signals: 0,
        }
    }

    fn push_button(&mut self) {
        self.presses += 1;
        self.queue.push_back(Signal {
            from: "button".to_string(),
            pulse: Pulse::Low,
            to: "broadcaster".to_string(),
        });
        // let mut rx_low = 0usize;
        // let mut rx_high = 0usize;
        while let Some(signal) = self.queue.pop_front() {
            // dbg!(&signal);
            if signal.pulse == Pulse::Low {
                self.low_signals += 1;
            } else {
                self.high_signals += 1;
            }
            // if signal.to == "rx" {
            //     if signal.pulse == Pulse::Low {
            //         rx_low += 1;
            //     } else {
            //         rx_high += 1;
            //     }
            // }
            if let Some(node) = self.nodes.get_mut(&signal.to) {
                node.receive(signal, &mut self.queue);
            }
        }
        // if rx_low == 1 && rx_high == 0 {
        //     Some(self.presses)
        // } else {
        //     None
        // }
    }

    fn total(&self) -> usize {
        self.low_signals * self.high_signals
    }
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let mut network = Network::parse(reader);
    for _ in 0..1000 {
        network.push_button();
    }
    println!("Total: {}", network.total());
    Ok(())
}
