import React from "react";
import { expect, test } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";

test("Title is rendered", async () => {
  render(<App />);

  const title = screen.getByText("Balance Sheet");

  expect(title).toBeInTheDocument();
});
