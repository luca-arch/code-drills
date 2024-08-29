import React from "react";
import { Breadcrumb, Typography } from "antd";
import HeaderContainer from "./HeaderContainer";
import SectionContainer from "./SectionContainer";
import type { Report } from "../api";

interface ComponentProps {
  report: Report;
}

const ReportContainer: React.FC<ComponentProps> = ({ report }) => {
  const breadcrumbItems = report.ReportTitles.map((title) => ({ title }));

  return (
    <>
      <Breadcrumb items={breadcrumbItems} />

      <Typography.Title level={3}>
        Report for {report.ReportName} (id: {report.ReportID})
      </Typography.Title>

      {report.Rows.map((row, key) => {
        switch (row.RowType) {
          case "Header":
            return <HeaderContainer key={key} row={row} />;
          case "Section":
            return <SectionContainer key={key} row={row} />;
          default:
            return <div>Error: this should never render!</div>;
        }
      })}
    </>
  );
};

export default ReportContainer;
