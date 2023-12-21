use core::panic;
use std::cmp::Ordering;
use std::collections::HashMap;
use std::fs::File;
use std::io::BufRead;

type Part = HashMap<char, usize>;

fn parse_part(line: &str) -> Part {
    let mut part = Part::new();
    let e = line.find('}').unwrap();
    let mut ch = line[..e].chars();
    assert_eq!(ch.next(), Some('{'));
    for kv in ch.as_str().split(',') {
        let mut kv = kv.split('=');
        let k = kv.next().unwrap().chars().next().unwrap();
        let v = kv.next().unwrap().parse::<usize>().unwrap();
        part.insert(k, v);
    }
    part
}

struct Rule {
    field: char,
    cmp: Ordering,
    value: usize,
    dest: String,
}

impl Rule {
    fn parse(line: &str) -> Rule {
        let mut parts = line.split(':');
        let mut predicate = parts.next().unwrap().trim().chars();
        let dest = parts.next().unwrap().trim().to_string();
        
        let field = predicate.next().unwrap();
        let cmp = match predicate.next().unwrap() {
            '<' => Ordering::Less,
            '>' => Ordering::Greater,
            '=' => Ordering::Equal,
            _ => panic!("Invalid comparison operator"),
        };
        Rule {
            field,
            cmp,
            value: predicate.as_str().parse().unwrap(),
            dest,
        }
    }

    fn test(&self, part: &Part) -> Option<&str> {
        let value = part.get(&self.field).unwrap();
        let cmp = value.cmp(&self.value);
        if cmp == self.cmp {
            Some(&self.dest)
        } else {
            None
        }
    }
}

struct RuleSet {
    rules: Vec<Rule>,
    default: String,
}

impl RuleSet {
    fn parse(line: &str) -> (String, RuleSet) {
        let i = line.find('{').unwrap();
        let name = line[..i].trim().to_string();
        let j = line.find('}').unwrap();
        let parts = line[i + 1..j].split(',').collect::<Vec<_>>();
        let k = parts.len() - 1;
        (name, RuleSet {
            rules: parts[..k]
                .iter()
                .map(|&s| Rule::parse(s)).collect(),
            default: parts[k].to_string(),
        })
    }

    fn test(&self, part: &Part) -> &str {
        for rule  in &self.rules {
            if let Some(dest) = rule.test(part) {
                return dest;
            }
        }
        &self.default
    }
}

fn parse_input(reader: impl BufRead) -> (HashMap<String, RuleSet>, Vec<Part>) {
    let mut rules = HashMap::new();
    let mut parts = Vec::new();

    let mut lines = reader.lines();
    for line in lines.by_ref() {
        let line = line.unwrap();
        if line.is_empty() {
            break;
        }
        let (name, rs) = RuleSet::parse(&line);
        rules.insert(name, rs);
    }
    for line in lines {
        let line =line.unwrap();
        parts.push(parse_part(&line));
    }
    (rules, parts)
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let (rules, parts) = parse_input(reader);

    let mut accepted = Vec::new();
    for part in &parts {
        print!("{:?}: ", part);
        let mut name = "in";
        loop {
            print!("{}", name);   
            let rs = rules.get(name).unwrap();
            match rs.test(part) {
                "A" => {
                    println!(" -> A");   
                    accepted.push(part);
                    break;
                },
                "R" => {
                    println!(" -> R");   
                    break;
                },
                x => {
                    print!(" -> ");   
                    name = x;
                },
            }
        }
    }
    
    let mut sum = 0;
    for part in accepted {
        for &v in part.values() {
            sum += v;
        }
    }
    println!("Sum: {}", sum);

    // TODO: Solve a problem!
    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_part_parse() {
        let part = parse_part("{a=1,b=2,c=3}");
        assert_eq!(part.get(&'a'), Some(&1));
        assert_eq!(part.get(&'b'), Some(&2));
        assert_eq!(part.get(&'c'), Some(&3));
    }

    #[test]
    fn test_rule_parse() {
        let rule = Rule::parse("a<2006:qkq");
        assert_eq!(rule.field, 'a');
        assert_eq!(rule.cmp, Ordering::Less);
        assert_eq!(rule.value, 2006);
        assert_eq!(rule.dest, "qkq");
    }
}
