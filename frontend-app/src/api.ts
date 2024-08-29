import { Dayjs } from "dayjs";

const GET_REPORTS_ENDPOINT = "/balance";

export type CellAttr = {
  ID: string;
  Value: string;
};

export type Cell = {
  Attributes?: CellAttr[];
  Value: string;
};

/**
 * Response schema for `GET_REPORTS_ENDPOINT`.
 *
 * @link https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
 */
export type ReportResponse = {
  Reports: Report[];
};

/**
 * Search form parameters.
 *
 * @link https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
 */
export type ReportSearchParams = {
  /** Date */
  date?: Dayjs;

  /** Set this to true to get cash transactions only */
  paymentsOnly?: boolean;

  /** The number of periods to compare (integer between 1 and 11) */
  periods?: number;

  /** If you set this parameter to "true" then no custom report layouts will be applied to response */
  standardLayout?: boolean;

  /** The period size to compare to */
  timeframe?: "MONTH" | "QUARTER" | "YEAR";

  /** The balance sheet will be filtered by this option if supplied. Note you cannot filter just by the TrackingCategory */
  trackingOptionID1?: string;

  /** If you want to filter by more than one tracking category option then you can specify a second option too. See the Balance Sheet report in Xero learn more about this behavior when filtering by tracking category options */
  trackingOptionID2?: string;
};

export type Report = {
  ReportDate: string;
  ReportID: string;
  ReportName: string;
  ReportType: string;
  ReportTitles: string[];
  Rows: Row[];
  UpdatedDateUTC: Date;
};

export type Row = {
  Cells: Cell[];
  RowType: "Header" | "Row" | "Section" | "SummaryRow";
  Rows?: Row[];
  Title: string;
};

/**
 * Fetch Balance reports from the backend
 *
 * @fixme needs to use a decent client (like axios).
 */
export const getReports = (
  params: ReportSearchParams,
): Promise<ReportResponse> => {
  const httpParams = new URLSearchParams();

  for (const [name, value] of Object.entries(params)) {
    if (value === undefined) {
      continue;
    }

    if (name === "date") {
      httpParams.append(name, (value as Dayjs).toISOString());
    } else {
      httpParams.append(name, value.toString());
    }
  }

  const url = `${GET_REPORTS_ENDPOINT}?${httpParams.toString()}`;

  return fetch(url).then((res) => {
    if (!res.ok) {
      throw new Error(res.statusText);
    }

    return res.json();
  });
};
