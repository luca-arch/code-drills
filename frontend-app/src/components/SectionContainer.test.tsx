import React from "react";
import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import SectionContainer from "./SectionContainer";

import type { Row } from "../api";

test("Renders the summary row with a bold font", async () => {
  const mockRow: Row = {
    Cells: [],
    RowType: "Section",
    Rows: [
      {
        Cells: [
          { Attributes: [], Value: "Cell 1" },
          { Attributes: [], Value: "Cell 2" },
          { Attributes: [], Value: "Cell 3" },
        ],
        RowType: "SummaryRow",
        Rows: [],
        Title: "",
      },
    ],
    Title: "",
  };

  render(<SectionContainer row={mockRow} />);

  expect(screen.getByText("Cell 1")).toContainHTML("<strong>Cell 1</strong>");

  expect(screen.getByText("Cell 2")).toContainHTML("<strong>Cell 2</strong>");

  expect(screen.getByText("Cell 3")).toContainHTML("<strong>Cell 3</strong>");
});

test("Renders other rows without bold font", async () => {
  const mockRow: Row = {
    Cells: [],
    RowType: "Section",
    Rows: [
      {
        Cells: [
          { Attributes: [], Value: "Cell 1" },
          { Attributes: [], Value: "Cell 2" },
          { Attributes: [], Value: "Cell 3" },
        ],
        RowType: "Row",
        Rows: [],
        Title: "",
      },
    ],
    Title: "",
  };

  render(<SectionContainer row={mockRow} />);

  expect(screen.getByText("Cell 1")).toContainHTML("Cell 1");

  expect(screen.getByText("Cell 2")).toContainHTML("Cell 2");

  expect(screen.getByText("Cell 3")).toContainHTML("Cell 3");
});

test("Renders the title in a divider", async () => {
  const mockRow: Row = {
    Cells: [],
    RowType: "Section",
    Rows: [],
    Title: "The Section Title",
  };

  render(<SectionContainer row={mockRow} />);

  expect(screen.getByText("The Section Title")).toBeInTheDocument();
});
