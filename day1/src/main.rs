use std::fs::File;
use std::io::{BufRead, BufReader};

fn translate_digits(input: &str) -> String {
    input
        //.replace("zero", "0")
        .replace("one", "1")
        .replace("two", "2")
        .replace("three", "3")
        .replace("four", "4")
        .replace("five", "5")
        .replace("six", "6")
        .replace("seven", "7")
        .replace("eight", "8")
        .replace("nine", "9")
}

fn leading_digit(input: &str) -> (Option<u32>, usize) {
    let c = input.chars().next().unwrap();
    if c.is_ascii_digit() {
        return (Some(c.to_digit(10).unwrap()), 1);
    }
    if input.starts_with("one") {
        return (Some(1), 3);
    }
    if input.starts_with("two") {
        return (Some(2), 3);
    }
    if input.starts_with("three") {
        return (Some(3), 5);
    }
    if input.starts_with("four") {
        return (Some(4), 4);
    }
    if input.starts_with("five") {
        return (Some(5), 4);
    }
    if input.starts_with("six") {
        return (Some(6), 3);
    }
    if input.starts_with("seven") {
        return (Some(7), 5);
    }
    if input.starts_with("eight") {
        return (Some(8), 5);
    }
    if input.starts_with("nine") {
        return (Some(9), 4);
    }
    (None, 1)
}

fn get_digits(input: &str) -> Vec<u32> {
    let mut out = Vec::new();
    let mut i = 0;
    while i < input.len() {
        let (digit, step) = leading_digit(&input[i..]);
        if let Some(d) = digit {
            out.push(d);
        }
        //i += step;
        i += 1;
    }
    out
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();

    let mut sum = 0;
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);
    for l in reader.lines() {
        let line = l.unwrap();
        let digits = get_digits(&line);
        let first = digits.first().unwrap();
        let last = digits.last().unwrap();
        let number = first * 10 + last;
        println!("{:?} --> ({}, {}) --> {}", &line, first, last, number);
        sum += number;
    }
    println!("SUM = {}", sum);
    Ok(())
}
