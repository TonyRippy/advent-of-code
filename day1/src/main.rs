use std::cmp::Ordering;
use std::fs::File;
use std::io::BufRead;

const DIGITS: [(&str, u32); 18] = [
    ("1", 1),
    ("2", 2),
    ("3", 3),
    ("4", 4),
    ("5", 5),
    ("6", 6),
    ("7", 7),
    ("8", 8),
    ("9", 9),
    ("one", 1),
    ("two", 2),
    ("three", 3),
    ("four", 4),
    ("five", 5),
    ("six", 6),
    ("seven", 7),
    ("eight", 8),
    ("nine", 9),
];

fn find_digit<F>(find_func: F, ordering: Ordering) -> Option<u32>
where
    F: Fn(&str) -> Option<usize>,
{
    DIGITS
        .iter()
        .filter_map(|&(s, d)| find_func(s).map(|i| (i, d)))
        .reduce(|(ai, ad), (bi, bd)| {
            if bi.cmp(&ai) == ordering {
                (bi, bd)
            } else {
                (ai, ad)
            }
        })
        .map(|(_, d)| d)
}

fn find_first_digit(input: &str) -> Option<u32> {
    find_digit(|s| input.find(s), Ordering::Less)
}

fn find_last_digit(input: &str) -> Option<u32> {
    find_digit(|s| input.rfind(s), Ordering::Greater)
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let mut sum = 0;
    for l in reader.lines() {
        let line = l.unwrap();
        let first = find_first_digit(&line).unwrap();
        let last = find_last_digit(&line).unwrap();
        let number = first * 10 + last;
        println!("{:?} --> ({}, {}) --> {}", &line, first, last, number);
        sum += number;
    }
    println!("SUM = {}", sum);
    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_first() {
        assert_eq!(find_first_digit("1"), Some(1));
        assert_eq!(find_first_digit("a2c"), Some(2));
    }

    #[test]
    fn test_last() {
        assert_eq!(find_last_digit("1"), Some(1));
        assert_eq!(find_last_digit("4qtfn"), Some(4));
    }
}
