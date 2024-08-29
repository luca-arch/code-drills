import React from "react";
import { expect, test, vi } from "vitest";
import { fireEvent, render } from "@testing-library/react";
import SearchForm from "./SearchForm";

test("Renders all fields", async () => {
  const mockSetReports = vi.fn();

  const form = render(<SearchForm setReports={mockSetReports} />);

  expect(form.getByText("Date")).toBeInTheDocument();
  expect(form.getByText("Periods")).toBeInTheDocument();
  expect(form.getByText("Timeframe")).toBeInTheDocument();
  expect(form.getByText("Standard layout")).toBeInTheDocument();
  expect(form.getByText("Only cash")).toBeInTheDocument();

  expect(form.getByText("Search")).toBeInTheDocument();
});

test("Search button correctly updates the state", async () => {
  const mockSetReports = vi.fn();

  const form = render(<SearchForm setReports={mockSetReports} />);

  const searchButton = form.getByText("Search");

  expect(searchButton).toBeInTheDocument();

  fireEvent.click(searchButton);

  // TODO: test mockSetReports is called with empty ReportSearchParams
});
