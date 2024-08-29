import React from "react";
import { Typography } from "antd";
import Main from "./components/Main";

const App: React.FC = () => (
  <>
    <Typography.Title style={{ textAlign: "center" }}>
      Balance Sheet
    </Typography.Title>

    <Main />
  </>
);

export default App;
