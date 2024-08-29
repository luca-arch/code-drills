import React from "react";
import { Typography } from "antd";
import type { Row } from "../api";

interface ContainerProps {
  row: Row;
}

/**
 * HeaderContainer component renders rows of type "Header".
 */
const HeaderContainer: React.FC<ContainerProps> = ({ row }) => (
  <Typography.Title level={4}>
    {row.Cells.filter((cell) => cell.Value && cell.Value.length)
      .reverse()
      .map((cell) => cell.Value)
      .join(" - ")}
  </Typography.Title>
);

export default HeaderContainer;
