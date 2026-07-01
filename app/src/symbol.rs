// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Symbol table for compactly representing repeated strings.
use std::collections::HashMap;

/// A compact handle to an interned string.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Symbol(pub usize);

/// Stores unique strings and assigns stable symbols to them.
#[derive(Debug, Default)]
pub struct SymbolTable {
    /// Strings in symbol order.
    strings: Vec<String>,

    /// Reverse lookup from string contents to the symbol already assigned.
    lookup_map: HashMap<String, Symbol>,
}

impl SymbolTable {
    /// Creates an empty string interner.
    pub fn new() -> Self {
        Self::default()
    }

    /// Returns the symbol for `text`, inserting it if needed.
    ///
    /// If `text` was interned before, the existing symbol is returned.
    /// Otherwise, `text` is stored once and assigned the next symbol.
    pub fn intern(&mut self, text: &str) -> Symbol {
        if let Some(&symbol) = self.lookup_map.get(text) {
            return symbol;
        }

        let owned = text.to_owned();
        let symbol = Symbol(self.strings.len());

        self.strings.push(owned.clone());
        self.lookup_map.insert(owned, symbol);

        symbol
    }

    /// Resolves a symbol back to the interned string.
    ///
    /// # Panics
    ///
    /// Panics if `symbol` was not created by this interner.
    pub fn resolve(&self, symbol: Symbol) -> &str {
        self.strings[symbol.0].as_str()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn same_string_gets_same_symbol() {
        let mut table = SymbolTable::new();

        let first = table.intern("alpha");
        let second = table.intern("alpha");

        assert_eq!(first, second);
        assert_eq!(table.resolve(first), "alpha");
    }

    #[test]
    fn different_strings_get_different_symbols() {
        let mut table = SymbolTable::new();

        let first = table.intern("alpha");
        let second = table.intern("beta");

        assert_ne!(first, second);
        assert_eq!(table.resolve(first), "alpha");
        assert_eq!(table.resolve(second), "beta");
    }
}
