// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Maps byte offsets to line and column locations.
use super::LineColumn;
use super::Position;

/// An index for translating byte offsets into line and column locations.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct LineMap {
    line_starts: Vec<Position>,
}

impl LineMap {
    /// Creates a line map for the given source text.
    pub fn new(source: &str) -> Self {
        let mut line_starts = vec![super::Position(0)];

        for (offset, byte) in source.bytes().enumerate() {
            if byte == b'\n' {
                line_starts.push(Position(offset + 1));
            }
        }

        Self { line_starts }
    }

    /// Returns the line and column location for the given position.
    pub fn line_column(&self, position: Position) -> LineColumn {
        let line = self
            .line_starts
            .partition_point(|line_start| *line_start <= position)
            .saturating_sub(1);

        let line_start = self.line_starts[line];

        LineColumn {
            line: line + 1,                        // Convert to one-based line number.
            column: position.0 - line_start.0 + 1, // Convert to one-based column number.
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn line_map_maps_start_of_empty_source() {
        let line_map = LineMap::new("");
        let location = line_map.line_column(Position(0));

        assert_eq!(location, LineColumn { line: 1, column: 1 });
    }

    #[test]
    fn line_map_maps_start_of_source() {
        let line_map = LineMap::new("abc");
        let location = line_map.line_column(Position(0));

        assert_eq!(location, LineColumn { line: 1, column: 1 });
    }

    #[test]
    fn line_map_maps_position_on_first_line() {
        let line_map = LineMap::new("abc");
        let location = line_map.line_column(Position(2));

        assert_eq!(location, LineColumn { line: 1, column: 3 });
    }

    #[test]
    fn line_map_maps_start_of_second_line() {
        let line_map = LineMap::new("abc\ndef");
        let location = line_map.line_column(Position(4));

        assert_eq!(location, LineColumn { line: 2, column: 1 });
    }

    #[test]
    fn line_map_maps_position_inside_second_line() {
        let line_map = LineMap::new("abc\ndef");
        let location = line_map.line_column(Position(6));

        assert_eq!(location, LineColumn { line: 2, column: 3 });
    }

    #[test]
    fn line_map_maps_empty_line() {
        let line_map = LineMap::new("abc\n\ndef");
        let location = line_map.line_column(Position(5));

        assert_eq!(location, LineColumn { line: 3, column: 1 });
    }

    #[test]
    fn line_map_maps_position_after_multibyte_character() {
        let line_map = LineMap::new("é\nx");
        let location = line_map.line_column(Position(3));

        assert_eq!(location, LineColumn { line: 2, column: 1 });
    }
}
