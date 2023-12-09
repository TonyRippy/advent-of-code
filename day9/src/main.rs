use std::fs::File;
use std::io::BufRead;

struct Diff<'a, I>
where
    I: Iterator<Item = i64>,
{
    input: &'a mut I,
    last: Option<i64>,
}

impl<'a, I> Diff<'a, I>
where
    I: Iterator<Item = i64>,
{
    fn new(input: &'a mut I) -> Self {
        let last = input.next();
        Self { input, last }
    }
}

impl<'a, I> Iterator for Diff<'a, I>
where
    I: Iterator<Item = i64>,
{
    type Item = i64;
    fn next(&mut self) -> Option<Self::Item> {
        let next = self.input.next()?;
        let last = self.last?;
        self.last = Some(next);
        Some(next - last)
    }
}

fn next_in_sequence(sequence: Vec<i64>) -> Option<i64> {
    let last = sequence.last()?.clone();
    let mut diffs = vec![];
    let mut zeros = true;

    for d in Diff::new(&mut sequence.into_iter()) {
        if d != 0 {
            zeros = false;
        }
        diffs.push(d);
    }
    if zeros {
        return Some(last);
    }
    let last_diff = next_in_sequence(diffs)?;
    Some(last + last_diff)
}

fn prev_in_sequence(sequence: Vec<i64>) -> Option<i64> {
    let first = *sequence.first()?;
    let mut diffs = vec![];
    let mut zeros = true;

    for d in Diff::new(&mut sequence.into_iter()) {
        if d != 0 {
            zeros = false;
        }
        diffs.push(d);
    }
    if zeros {
        return Some(first);
    }
    let first_diff = prev_in_sequence(diffs)?;
    Some(first - first_diff)
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    // PART 1

    // let sum: i64 = reader
    //     .lines()
    //     .map(|l| {
    //         l.unwrap()
    //             .split_ascii_whitespace()
    //             .map(|n| n.parse::<i64>().unwrap())
    //             .collect::<Vec<_>>()
    //     })
    //     .map(|v| next_in_sequence(v).unwrap())
    //     .sum();

    // PART 2

    let sum: i64 = reader
        .lines()
        .map(|l| {
            l.unwrap()
                .split_ascii_whitespace()
                .map(|n| n.parse::<i64>().unwrap())
                .collect::<Vec<_>>()
        })
        .map(|v| prev_in_sequence(v).unwrap())
        .sum();

    println!("Sum: {}", sum);

    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_diffs() {
        let mut input = vec![1, 2, 3, 4, 5].into_iter();
        let mut diffs = Diff::new(&mut input);
        assert_eq!(diffs.next(), Some(1));
        assert_eq!(diffs.next(), Some(1));
        assert_eq!(diffs.next(), Some(1));
        assert_eq!(diffs.next(), Some(1));
        assert_eq!(diffs.next(), None);
    }

    #[test]
    fn test_zeros() {
        let input = vec![0, 0, 0];
        assert_eq!(next_in_sequence(input), Some(0));
    }

    #[test]
    fn test_ones() {
        let input = vec![1, 1, 1];
        assert_eq!(next_in_sequence(input), Some(1));
    }

    #[test]
    fn test_sequence() {
        let input = vec![1, 2, 3];
        assert_eq!(next_in_sequence(input), Some(4));
    }
}
