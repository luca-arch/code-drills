import React, { useState } from "react";

import { Col, Divider, Row } from "antd";
import ReportContainer from "./ReportContainer";
import SearchForm from "./SearchForm";
import type { Report } from "../api";

const Main: React.FC = () => {
  const [reports, setReports] = useState<Report[]>();

  return (
    <Row>
      <Col span={10}>
        <SearchForm setReports={setReports} />
      </Col>

      <Col span={14}>
        {reports &&
          reports.map((report) => (
            <>
              <Divider />

              <Col span={24}>
                <ReportContainer report={report} />
              </Col>
            </>
          ))}
      </Col>
    </Row>
  );
};

export default Main;
