export type ParsedLine = {
  type: "loguru" | "table" | "plain";
  level?: string;
  timestamp?: string;
  message: string;
  raw: string;
};

const LOGURU_RE =
  /^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})\s*\|\s*(\w+)\s*\|\s*\S+\s*-\s*(.*)$/;

export function parseAmdLine(text: string, _source: "stdout" | "stderr"): ParsedLine {
  const loguruMatch = text.match(LOGURU_RE);
  if (loguruMatch) {
    return {
      type: "loguru",
      timestamp: loguruMatch[1],
      level: loguruMatch[2],
      message: loguruMatch[3],
      raw: text,
    };
  }

  if (text.includes("+---") || (text.includes("|") && text.split("|").length > 2)) {
    return { type: "table", message: text, raw: text };
  }

  return { type: "plain", message: text, raw: text };
}

export function levelColor(level?: string): string {
  switch (level?.toUpperCase()) {
    case "SUCCESS":
      return "text-green-500";
    case "WARNING":
      return "text-yellow-500";
    case "ERROR":
    case "CRITICAL":
      return "text-red-500";
    default:
      return "text-text";
  }
}
