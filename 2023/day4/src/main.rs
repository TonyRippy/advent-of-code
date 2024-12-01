use std::collections::HashSet;
use std::fs::File;
use std::io::BufRead;

fn matches_in_card(line: &str) -> usize {
    let start = line.find(':').unwrap() + 1;
    let mid = line.find('|').unwrap();
    let winning = line[start..mid]
        .split_whitespace()
        .filter(|&x| !x.is_empty())
        .flat_map(|x| x.parse::<usize>())
        .collect::<HashSet<usize>>();
    let have = line[mid + 1..]
        .split_whitespace()
        .filter(|&x| !x.is_empty())
        .flat_map(|x| x.parse::<usize>())
        .collect::<HashSet<usize>>();
    winning.intersection(&have).count()
}

#[inline]
fn matches_to_points(matches: usize) -> usize {
    (1 << matches) >> 1
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    // PART 1

    // let sum = reader.lines().flatten()
    //     .map(|line| matches_in_card(&line))
    //     .map(matches_to_points)
    //     .sum::<usize>();

    // PART 2

    let lines = reader.lines().flatten().collect::<Vec<_>>();
    let n = lines.len();
    let mut counts = vec![1_usize; n];

    for i in 0..n {
        let count = counts[i];
        let matches = matches_in_card(&lines[i]);
        println!("{}: count = {}, matches = {}", i, count, matches);
        if matches > 0 {
            for j in i + 1..i + matches + 1 {
                counts[j] += count;
            }
        }
    }
    let sum = counts.iter().sum::<usize>();

    println!("SUM = {}", sum);

    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_matches_to_points() {
        assert_eq!(matches_to_points(0), 0);
        assert_eq!(matches_to_points(1), 1);
        assert_eq!(matches_to_points(2), 2);
        assert_eq!(matches_to_points(3), 4);
        assert_eq!(matches_to_points(4), 8);
    }
}
