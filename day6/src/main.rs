use std::fs::File;
use std::io::BufRead;

fn read_line_part_1(line: &str) -> (String, Vec<usize>) {
    let mut parts = line.split(':');
    let name = parts.next().unwrap().to_string();
    let values = parts
        .next()
        .unwrap()
        .split_whitespace()
        .map(|x| x.parse::<usize>().unwrap())
        .collect();
    (name, values)
}

fn read_line_part_2(line: &str) -> (String, usize) {
    let mut parts = line.split(':');
    let name = parts.next().unwrap().to_string();
    let mut value = String::default();
    parts
        .next()
        .unwrap()
        .split_whitespace()
        .for_each(|s| value.push_str(s));
    (name, value.parse::<usize>().unwrap())
}

fn ways_to_win(time: usize, distance: usize) -> usize {
    (0..time)
        .filter_map(|press| {
            let d = press * (time - press);
            if d > distance {
                Some((press, d))
            } else {
                None
            }
        })
        .count()
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    // PART 1

    // let mut it = reader.lines().map(|x| x.unwrap());
    // let (name, time) = read_line_part_1(it.next().unwrap().as_str());
    // assert_eq!(name, "Time");
    // let (name, distance) = read_line_part_1(it.next().unwrap().as_str());
    // assert_eq!(name, "Distance");
    // assert!(it.next().is_none());

    // let answer: usize = time
    //     .iter()
    //     .zip(distance.iter())
    //     .map(|(t, d)| {
    //         print!("Time: {}, Distance: {}", t, d);
    //         let ways = ways_to_win(*t, *d);
    //         println!(" Ways to win: {}", ways);
    //         ways
    //     })
    //     .product();

    // println!("Answer: {}", answer);

    // PART 2

    let mut it = reader.lines().map(|x| x.unwrap());
    let (name, time) = read_line_part_2(it.next().unwrap().as_str());
    assert_eq!(name, "Time");
    let (name, distance) = read_line_part_2(it.next().unwrap().as_str());
    assert_eq!(name, "Distance");
    assert!(it.next().is_none());
    
    print!("Time: {}, Distance: {}", time, distance);
    let ways = ways_to_win(time, distance);
    println!(" Ways to win: {}", ways);

    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_something() {}
}
