import React from "react";
import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import HeaderContainer from "./HeaderContainer";

import type { Row } from "../api";

test("Renders the headers cells in reverse order", async () => {
  const mockRow: Row = {
    Cells: [
      { Attributes: [], Value: "A" },
      { Attributes: [], Value: "B" },
      { Attributes: [], Value: "C" },
    ],
    RowType: "Header",
    Rows: [],
    Title: "",
  };

  render(<HeaderContainer row={mockRow} />);

  const elem = screen.getByText("C - B - A");

  expect(elem).toBeInTheDocument();
});

test("Discards empty cells", async () => {
  const mockRow: Row = {
    Cells: [
      { Attributes: [], Value: "" },
      { Attributes: [], Value: "A" },
      { Attributes: [], Value: "" },
      { Attributes: [], Value: "C" },
      { Attributes: [], Value: "" },
    ],
    RowType: "Header",
    Rows: [],
    Title: "",
  };

  render(<HeaderContainer row={mockRow} />);

  const elem = screen.getByText("C - A");

  expect(elem).toBeInTheDocument();
});
