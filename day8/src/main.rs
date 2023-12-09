use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::BufRead;

fn lcm(a: usize, b: usize) -> usize {
    a * b / gcd(a, b)
}

fn gcd(a: usize, b: usize) -> usize {
    let mut a = a;
    let mut b = b;
    while b != 0 {
        let t = b;
        b = a % b;
        a = t;
    }
    a
}

fn distance_between_nodes(
    from: &str,
    to: &str,
    directions: &str,
    graph: &HashMap<String, (String, String)>,
) -> Option<usize> {
    let dir = directions.chars().collect::<Vec<_>>();
    let mut d_idx = 0;

    // Keep track of where we've been so we don't get stuck in a loop
    let mut visited = HashSet::new();

    let mut steps = 0;
    let mut pos = from;
    while pos != to {
        if visited.contains(&(pos, d_idx)) {
            return None;
        }
        visited.insert((pos, d_idx));

        let (left, right) = &graph[pos];
        let next = match dir[d_idx] {
            'L' => left,
            'R' => right,
            _ => panic!("Unknown direction"),
        };
        d_idx = (d_idx + 1) % dir.len();
        pos = next;
        steps += 1;
    }
    Some(steps)
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    // Parse the input
    let mut lines = reader.lines().flatten();
    let directions = lines.next().unwrap();
    assert!(lines.next().unwrap().is_empty()); // skip blank line
    let graph = lines
        .map(|line| {
            let mut line = line.split(" = ");
            let from = line.next().unwrap().trim();
            let mut to = line.next().unwrap().trim().split(", ");
            let left = to.next().unwrap();
            let right = to.next().unwrap();
            assert!(to.next().is_none());
            assert_eq!(left.chars().next().unwrap(), '(');
            assert_eq!(right.chars().last().unwrap(), ')');
            (
                from.to_string(),
                (left[1..].to_string(), right[..right.len() - 1].to_string()),
            )
        })
        .collect::<HashMap<String, (String, String)>>();
    // println!("Graph = {:?}", graph);

    // PART 1

    // // Walk the map
    // let steps = distance_between_nodes("AAA", "ZZZ", &directions, &graph);
    // println!("Steps = {}", steps);

    // PART 2

    // find starting positions
    let starting_nodes = graph
        .keys()
        .filter(|&k| k.ends_with('A'))
        .collect::<Vec<_>>();
    let ending_nodes = graph
        .keys()
        .filter(|&k| k.ends_with('Z'))
        .collect::<Vec<_>>();
    let mut all_distances = vec![];
    for n in starting_nodes.iter() {
        let distances = ending_nodes
            .iter()
            .flat_map(|m| distance_between_nodes(n, m, &directions, &graph))
            .collect::<Vec<_>>();
        println!("{} = {:?}", n, &distances);
        // This was not obvious to me, but I guess each starting node can only reach one ending node???
        // That certainly makes it easier to calculate a solution.
        assert_eq!(distances.len(), 1);
        all_distances.push(distances[0]);
    }
    println!("Distances = {:?}", &all_distances);
    println!("LCM = {:?}", all_distances.into_iter().reduce(lcm));

    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_something() {}
}
