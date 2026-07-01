// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Loaded source files.
use std::path::{Path, PathBuf};

use super::FileId;

use crate::location::LineMap;

/// A source file loaded for analysis.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SourceFile {
    id: FileId,
    path: PathBuf,
    text: String,
    line_map: LineMap,
}

impl SourceFile {
    /// Creates a source file.
    pub fn new(id: FileId, path: PathBuf, text: String) -> Self {
        let line_map = LineMap::new(&text);

        Self {
            id,
            path,
            text,
            line_map,
        }
    }

    /// Returns the source file identifier.
    pub const fn id(&self) -> FileId {
        self.id
    }

    /// Returns the source file path.
    pub fn path(&self) -> &Path {
        &self.path
    }

    /// Returns the source text.
    pub fn text(&self) -> &str {
        &self.text
    }

    /// Returns the line map for the source text.
    pub const fn line_map(&self) -> &LineMap {
        &self.line_map
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::location::{LineColumn, Position};

    #[test]
    fn source_file_preserves_id_path_and_text() {
        let source_file = SourceFile::new(
            FileId(7),
            PathBuf::from("src/main.rSpot"),
            String::from("abc"),
        );

        assert_eq!(source_file.id(), FileId(7));
        assert_eq!(source_file.path(), PathBuf::from("src/main.rSpot"));
        assert_eq!(source_file.text(), "abc");
    }

    #[test]
    fn source_file_builds_line_map_from_text() {
        let source_file = SourceFile::new(
            FileId(1),
            PathBuf::from("src/main.rSpot"),
            String::from("abc\ndef"),
        );

        let location = source_file.line_map().line_column(Position(4));

        assert_eq!(location, LineColumn { line: 2, column: 1 });
    }
}
