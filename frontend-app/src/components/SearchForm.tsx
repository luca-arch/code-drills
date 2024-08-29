// @ts-expect-error TS6133
import React, { useState } from "react";

import { Button, Checkbox, DatePicker, Form, InputNumber, Select } from "antd";
import { getReports, ReportSearchParams } from "../api";
import type { FormProps } from "antd";
import type { Report } from "../api";

interface ComponentProps {
  setReports: (report: Report[]) => void;
}

/**
 * The Balance search form component.
 */
const SearchForm = ({ setReports }: ComponentProps) => {
  const [errors, setErrors] = useState<string[]>();

  const onFinish: FormProps<ReportSearchParams>["onFinish"] = (params) => {
    getReports(params)
      .then((res) => {
        setReports(res.Reports);
        setErrors([]);
      })
      .catch((e: Error) => setErrors([e.message]));
  };

  const onFinishFailed: FormProps<ReportSearchParams>["onFinishFailed"] = (
    errorInfo,
  ) => {
    console.table(errorInfo);
  };

  return (
    <Form
      autoComplete="off"
      initialValues={{ remember: true }}
      labelCol={{ span: 8 }}
      name="basic"
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      style={{ maxWidth: 600 }}
      wrapperCol={{ span: 16 }}
    >
      <Form.Item<ReportSearchParams> label="Date" name="date">
        <DatePicker />
      </Form.Item>

      <Form.Item<ReportSearchParams>
        help="Number of periods to compare"
        label="Periods"
        name="periods"
      >
        <InputNumber max={11} min={1} />
      </Form.Item>

      <Form.Item<ReportSearchParams>
        help="Period size to compare to"
        label="Timeframe"
        name="timeframe"
      >
        <Select
          options={[
            { value: "MONTH", label: <span>Month</span> },
            { value: "QUARTER", label: <span>Quarter</span> },
            { value: "YEAR", label: <span>Year</span> },
          ]}
        />
      </Form.Item>

      <Form.Item<ReportSearchParams>
        label="Standard layout"
        name="standardLayout"
        valuePropName="checked"
      >
        <Checkbox>Do not show custom report layouts</Checkbox>
      </Form.Item>

      <Form.Item<ReportSearchParams>
        label="Only cash"
        name="paymentsOnly"
        valuePropName="checked"
      >
        <Checkbox>Show cash transactions only</Checkbox>
      </Form.Item>

      <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
        <Button type="primary" htmlType="submit">
          Search
        </Button>
      </Form.Item>

      {errors && errors.length ? (
        <>
          <strong>Error occurred!</strong>

          <Form.ErrorList
            errors={errors.map((errMsg, index) => (
              <span key={index}>{errMsg}</span>
            ))}
          />
        </>
      ) : null}
    </Form>
  );
};

export default SearchForm;
