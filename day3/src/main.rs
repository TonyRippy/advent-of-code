use std::fs::File;
use std::io::BufRead;

#[derive(Debug)]
struct Number {
    x1: usize,
    x2: usize,
    y: usize,
}

impl Number {
    fn adjacent_x_range(&self) -> (usize, usize) {
        let min_x = if self.x1 == 0 { 0 } else { self.x1 - 1 };
        let max_x = self.x2 + 1;
        assert!(min_x <= max_x);
        (min_x, max_x)
    }

    fn adjacent_y_range(&self) -> (usize, usize) {
        let min_y = if self.y == 0 { 0 } else { self.y - 1 };
        let max_y = self.y + 1;
        (min_y, max_y)
    }

    pub fn is_next_to_any_symbol(&self, symbols: &Vec<Vec<usize>>) -> bool {
        let (min_x, max_x) = self.adjacent_x_range();
        let (min_y, mut max_y) = self.adjacent_y_range();
        if max_y >= symbols.len() {
            max_y = symbols.len() - 1;
        }
        for s in &symbols[min_y..=max_y] {
            for &x in s {
                if x >= min_x && x <= max_x {
                    return true;
                }
            }
        }
        false
    }

    pub fn is_next_to(&self, x: usize, y: usize) -> bool {
        let (min_y, max_y) = self.adjacent_y_range();
        if y < min_y || y > max_y {
            return false;
        }
        let (min_x, max_x) = self.adjacent_x_range();
        x >= min_x && x <= max_x
    }

    pub fn get_value(&self, lines: &[String]) -> usize {
        let line = &lines[self.y];
        line[self.x1..=self.x2].parse::<usize>().unwrap()
    }
}

fn find_numbers_in_line(line_number: usize, line: &str) -> Vec<Number> {
    let digits = line
        .chars()
        .enumerate()
        .map(|(i, c)| (i, c.is_ascii_digit()))
        .collect::<Vec<_>>();
    digits
        .split(|(_, is_digit)| !is_digit)
        .filter(|d| !d.is_empty())
        .map(|d| {
            let start = d[0].0;
            let end = d[d.len() - 1].0;
            Number {
                x1: start,
                x2: end,
                y: line_number,
            }
        })
        .collect::<Vec<_>>()
}

fn find_symbols_in_line(line: &str) -> Vec<usize> {
    line.chars()
        .enumerate()
        .filter(|&(_, c)| !(c == '.' || c.is_ascii_digit()))
        .map(|(i, _)| i)
        .collect::<Vec<_>>()
}

fn find_all_gears(lines: &[String]) -> Vec<(usize, usize)> {
    lines
        .iter()
        .enumerate()
        .flat_map(|(y, l)| {
            let mut gears = Vec::new();
            let mut start = 0;
            while let Some(x) = l[start..].find('*') {
                let xx = start + x;
                gears.push((xx, y));
                start = xx + 1;
            }
            gears
        })
        .collect::<Vec<_>>()
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    // let fname = "test.txt";
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);
    let lines = reader.lines().map(|l| l.unwrap()).collect::<Vec<_>>();

    let numbers = lines
        .iter()
        .enumerate()
        .map(|(i, line)| find_numbers_in_line(i, line))
        .collect::<Vec<_>>();
    println!("Numbers = {:?}", &numbers);

    let mut sum = 0;

    // PART 1

    // let symbols = lines.iter().map(|line| find_symbols_in_line(line)).collect::<Vec<_>>();
    // println!("Symbols = {:?}", &symbols);
    //
    // for n in numbers.iter().flatten() {
    //     print!("Number = {:?}", n);
    //     if n.is_next_to_any_symbol(&symbols) {
    //         print!(" - is next to symbol");
    //         let value = n.get_value(&lines);
    //         println!(" (value = {})", value);
    //         sum += value;
    //     } else {
    //         println!(" - not next to symbol");
    //     }
    // }

    // PART 2

    let gears = find_all_gears(&lines);
    println!("Gears = {:?}", &gears);

    for gear in gears {
        let next_to_gear = numbers
            .iter()
            .flatten()
            .filter(|n| n.is_next_to(gear.0, gear.1))
            .collect::<Vec<_>>();
        println!("Gear = {:?} - next to {:?}", gear, next_to_gear);
        if next_to_gear.len() != 2 {
            continue;
        }
        sum += next_to_gear[0].get_value(&lines) * next_to_gear[1].get_value(&lines);
    }

    println!("Sum = {}", sum);
    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_something() {}
}
