import React from "react";
import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import ReportContainer from "./ReportContainer";

import type { Report } from "../api";

const mockReport: Report = {
  ReportDate: "23 February 2018",
  ReportID: "m0ck-report-id",
  ReportName: "Test Balance",
  ReportType: "BalanceSheet",
  ReportTitles: [
    "Title - Balance Sheet",
    "Title - Demo Company (AU)",
    "Title - As at 28 February 2018",
  ],
  Rows: [],
  UpdatedDateUTC: new Date(),
};

test("Renders the heading", async () => {
  render(<ReportContainer report={mockReport} />);

  expect(
    screen.getByText("Report for Test Balance (id: m0ck-report-id)"),
  ).toBeInTheDocument();
});

test("Renders the title breadcrumbs", async () => {
  render(<ReportContainer report={mockReport} />);

  expect(screen.getByText("Title - Balance Sheet")).toBeInTheDocument();
  expect(screen.getByText("Title - Demo Company (AU)")).toBeInTheDocument();
  expect(
    screen.getByText("Title - As at 28 February 2018"),
  ).toBeInTheDocument();
});

test("Renders the rows", async () => {
  const mockReport2 = Object.assign({}, mockReport, {
    Rows: [
      {
        Cells: [
          { Attributes: [], Value: "Header" },
          { Attributes: [], Value: "Heading" },
          { Attributes: [], Value: "Mock" },
        ],
        RowType: "Header",
        Rows: [],
        Title: "The Mock Heading Header",
      },
      {
        Cells: [],
        RowType: "Section",
        Rows: [
          {
            Cells: [
              { Attributes: [], Value: "Sumary Cell 1" },
              { Attributes: [], Value: "Sumary Cell 2" },
            ],
            RowType: "SummaryRow",
            Rows: [],
            Title: "",
          },
        ],
        Title: "The Mock Heading Section",
      },
      {
        Cells: [],
        RowType: "Section",
        Rows: [], // Should render "No items to show"
        Title: "Another Mock Heading Section",
      },
      {
        Cells: [],
        RowType: "Section",
        Rows: [
          {
            Cells: [
              { Attributes: [], Value: "Child Row 1, Cell 1" },
              { Attributes: [], Value: "Child Row 1, Cell 2" },
              { Attributes: [], Value: "Child Row 1, Cell 3" },
            ],
            RowType: "Row",
            Rows: [],
            Title: "",
          },
          {
            Cells: [
              { Attributes: [], Value: "Child Row 2, Cell 1" },
              { Attributes: [], Value: "Child Row 2, Cell 2" },
              { Attributes: [], Value: "Child Row 2, Cell 3" },
            ],
            RowType: "Row",
            Rows: [],
            Title: "",
          },
          {
            Cells: [
              { Attributes: [], Value: "Child Row 3, Cell 1" },
              { Attributes: [], Value: "Child Row 3, Cell 2" },
              { Attributes: [], Value: "Child Row 3, Cell 3" },
            ],
            RowType: "Row",
            Rows: [],
            Title: "",
          },
        ],
        Title: "",
      },
    ],
  });

  render(<ReportContainer report={mockReport2} />);

  expect(screen.getByText("Mock - Heading - Header")).toBeInTheDocument();

  expect(screen.getByText("The Mock Heading Section")).toBeInTheDocument();

  expect(screen.getByText("Sumary Cell 1")).toBeInTheDocument();

  expect(screen.getByText("Sumary Cell 2")).toBeInTheDocument();

  expect(screen.getByText("No items to show")).toBeInTheDocument();

  expect(screen.getByText("Another Mock Heading Section")).toBeInTheDocument();

  const expectedChildren = [
    "Child Row 1, Cell 1",
    "Child Row 1, Cell 2",
    "Child Row 1, Cell 3",
    "Child Row 2, Cell 1",
    "Child Row 2, Cell 2",
    "Child Row 2, Cell 3",
    "Child Row 3, Cell 1",
    "Child Row 3, Cell 2",
    "Child Row 3, Cell 3",
  ];

  expectedChildren.forEach((expectedText) => {
    expect(screen.getByText(expectedText)).toBeInTheDocument();
  });
});
