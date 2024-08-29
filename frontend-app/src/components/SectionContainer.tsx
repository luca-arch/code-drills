import React from "react";
import { Col, Divider, Row } from "antd";
import type { Cell, Row as ReportRow } from "../api";

interface InnerRowProps {
  cells: Cell[];
  isSummary?: boolean;
}

/**
 * Renders Section's children rows.
 * Based on stubbed response, this assumes the inner rows never contain more than 3 cells so it renders columns with span = 8/24
 */
const InnerRow: React.FC<InnerRowProps> = ({ cells, isSummary }) => (
  <Row>
    {cells &&
      cells.length &&
      cells.map((cell, cellNumber) => (
        <Col span={8} key={cellNumber}>
          {isSummary ? <strong>{cell.Value}</strong> : cell.Value}
        </Col>
      ))}
  </Row>
);

interface SectionContainerProps {
  row: ReportRow;
}

/**
 * SectionContainer component.
 */
const SectionContainer: React.FC<SectionContainerProps> = ({ row }) => (
  <>
    <Divider orientation="left">{row.Title}</Divider>

    {row.Rows && row.Rows.length ? (
      row.Rows.map((innerRow, innerRowNumber) => (
        <InnerRow
          key={innerRowNumber}
          cells={innerRow.Cells}
          isSummary={innerRow.RowType === "SummaryRow"}
        />
      ))
    ) : (
      <>No items to show</>
    )}
  </>
);

export default SectionContainer;
