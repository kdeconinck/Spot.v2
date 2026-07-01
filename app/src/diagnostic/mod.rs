// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Diagnostic types and rendering for rSpot.
mod diagnostic;
mod render;
mod severity;

pub use diagnostic::Diagnostic;
pub use render::render;
pub use severity::Severity;
