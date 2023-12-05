use std::fs::File;
use std::io::BufRead;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Range {
    start: usize,
    len: usize,
}

#[derive(Debug)]
struct MapEntry {
    dst_start: usize,
    src_start: usize,
    length: usize,
}

impl MapEntry {
    fn parse(line: &str) -> Self {
        let parts = line
            .split_ascii_whitespace()
            .flat_map(|part| part.parse::<usize>())
            .collect::<Vec<usize>>();
        Self {
            dst_start: parts[0],
            src_start: parts[1],
            length: parts[2],
        }
    }
}

#[derive(Debug)]
struct Map {
    name: String,
    entries: Vec<MapEntry>,
}

impl Map {
    fn parse<I: Iterator<Item = String>>(lines: &mut I) -> Option<Self> {
        // Parse name
        let name: String;
        loop {
            match lines.next() {
                Some(line) => {
                    // Skip any empty lines at the start
                    if line.is_empty() {
                        continue;
                    }
                    let i = line.find(" map:").unwrap();
                    name = line[..i].to_string();
                    break;
                }
                None => return None,
            }
        }

        // Parse entries
        let mut entries = Vec::new();
        for line in lines {
            if line.is_empty() {
                break;
            }
            entries.push(MapEntry::parse(&line))
        }
        entries.sort_by_key(|entry| entry.src_start);
        Some(Self { name, entries })
    }

    fn find(&self, src: usize) -> usize {
        for entry in &self.entries {
            if src >= entry.src_start {
                let delta = src - entry.src_start;
                if delta < entry.length {
                    let dst = entry.dst_start + delta;
                    return dst;
                }
            }
        }
        src
    }

    fn find_range(&self, mut src: Range) -> Vec<Range> {
        let mut out = Vec::new();
        for entry in self.entries.iter() {
            if src.start < entry.src_start {
                let delta = entry.src_start - src.start;
                if delta >= src.len {
                    out.push(src);
                    return out;
                }
                out.push(Range {
                    start: src.start,
                    len: delta,
                });
                src.start += delta;
                src.len -= delta;
            }
            if src.start >= entry.src_start {
                let offset = src.start - entry.src_start;
                if offset > entry.length {
                    continue;
                }
                let delta = src.len.min(entry.length - offset);
                out.push(Range {
                    start: entry.dst_start + offset,
                    len: delta,
                });
                src.start += delta;
                src.len -= delta;
            }
            if src.len == 0 {
                return out;
            }
        }
        out.push(src);
        out
    }
}

fn parse_input<I: Iterator<Item = String>>(lines: &mut I) -> (Vec<usize>, Vec<Map>) {
    // Parse seeds
    let line = lines.next().unwrap();
    let i = line.find(':').unwrap();
    let seeds = line[i + 2..]
        .split_ascii_whitespace()
        .flat_map(|part| part.parse::<usize>())
        .collect::<Vec<usize>>();

    // Parse maps
    let mut maps = Vec::new();
    while let Some(map) = Map::parse(lines) {
        maps.push(map);
    }
    (seeds, maps)
}

fn seeds_as_ranges(seeds: &[usize]) -> Vec<Range> {
    (0..seeds.len())
        .step_by(2)
        .map(|i| Range {
            start: seeds[i],
            len: seeds[i + 1],
        })
        .collect()
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let (seeds, maps) = parse_input(&mut reader.lines().map_while(Result::ok));

    // PART 1

    // let nearest = seeds
    //     .iter()
    //     .map(|&seed| {
    //         let mut dst = seed;
    //         print!("seed {}", &seed);
    //         for map in &maps {
    //             dst = map.find(dst);
    //             print!(" -> {}", dst);
    //         }
    //         println!();
    //         (seed, dst)
    //     })
    //     .min_by_key(|&(_, location)| location)
    //     .unwrap();

    // PART 2

    let nearest = seeds_as_ranges(&seeds)
        .into_iter()
        .flat_map(|seed_range| {
            let mut dsts = vec![seed_range];
            for map in &maps {
                dsts = dsts
                    .into_iter()
                    .flat_map(|dst| map.find_range(dst))
                    .collect();
            }
            dsts
        })
        .min_by_key(|&range| range.start)
        .unwrap();

    println!("Nearest: {:?}", nearest);

    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn find_before() {
        let map = Map {
            name: "test".to_string(),
            entries: vec![MapEntry {
                dst_start: 100,
                src_start: 10,
                length: 10,
            }],
        };
        let result = map.find_range(Range { start: 5, len: 2 });
        assert_eq!(result, vec![Range { start: 5, len: 2 }]);
    }

    #[test]
    fn find_before_and_in() {
        let map = Map {
            name: "test".to_string(),
            entries: vec![MapEntry {
                dst_start: 100,
                src_start: 10,
                length: 10,
            }],
        };
        let result = map.find_range(Range { start: 5, len: 10 });
        assert_eq!(
            result,
            vec![Range { start: 5, len: 5 }, Range { start: 100, len: 5 },]
        );
    }

    #[test]
    fn find_before_and_after() {
        let map = Map {
            name: "test".to_string(),
            entries: vec![MapEntry {
                dst_start: 100,
                src_start: 10,
                length: 10,
            }],
        };
        let result = map.find_range(Range { start: 5, len: 20 });
        assert_eq!(
            result,
            vec![
                Range { start: 5, len: 5 },
                Range {
                    start: 100,
                    len: 10
                },
                Range { start: 20, len: 5 }
            ]
        );
    }
}
