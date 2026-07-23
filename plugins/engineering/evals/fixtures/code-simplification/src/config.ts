export type Value = string | number | boolean;
export type Parsed = { values: Record<string, Value>; warnings: string[] };

export function parseConfig(
  source: string,
  env: Record<string, string | undefined>,
  onWarning: (message: string) => void,
): Parsed {
  const values: Record<string, Value> = {};
  const warnings: string[] = [];
  const lines = source.replace(/\r\n/g, "\n").split("\n");

  for (let index = 0; index < lines.length; index += 1) {
    const raw = lines[index];
    const lineNumber = index + 1;
    const trimmed = raw.trim();

    if (trimmed.length === 0) {
      continue;
    }
    if (trimmed.startsWith("#")) {
      continue;
    }

    const separator = raw.indexOf("=");
    if (separator < 1) {
      throw new SyntaxError(`Invalid config at line ${lineNumber}`);
    }

    const key = raw.slice(0, separator).trim();
    let value = raw.slice(separator + 1).trim();
    if (!/^[A-Za-z][A-Za-z0-9_.-]*$/.test(key)) {
      throw new SyntaxError(`Invalid key "${key}" at line ${lineNumber}`);
    }

    if (Object.prototype.hasOwnProperty.call(values, key)) {
      const warning = `Duplicate key "${key}" at line ${lineNumber}; last value wins`;
      warnings.push(warning);
      onWarning(warning);
    }

    value = value.replace(
      /\$\{([A-Za-z_][A-Za-z0-9_]*)\}/g,
      (_match, name: string) => {
        const replacement = env[name];
        if (replacement === undefined) {
          throw new ReferenceError(
            `Missing environment variable "${name}" at line ${lineNumber}`,
          );
        }
        return replacement;
      },
    );

    if (value === "true") {
      values[key] = true;
    } else if (value === "false") {
      values[key] = false;
    } else if (/^-?(0|[1-9][0-9]*)(\.[0-9]+)?$/.test(value)) {
      const parsedNumber = Number(value);
      if (!Number.isFinite(parsedNumber)) {
        throw new RangeError(
          `Non-finite number for "${key}" at line ${lineNumber}`,
        );
      }
      values[key] = parsedNumber;
    } else if (value.startsWith('"') && value.endsWith('"')) {
      values[key] = value.slice(1, -1);
    } else if (value.startsWith('"') || value.endsWith('"')) {
      throw new SyntaxError(
        `Unbalanced quote for "${key}" at line ${lineNumber}`,
      );
    } else {
      values[key] = value;
    }
  }

  return { values, warnings };
}
