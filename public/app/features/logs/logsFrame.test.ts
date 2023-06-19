import {
  DataFrame,
  FieldType,
  DataFrameType,
  dateTimeFormatISO,
  DateTimeInput,
  DateTimeOptions,
  Field,
  Labels,
} from '@grafana/data';

import { parseLogsFrame } from './logsFrame';

function makeString(name: string, values: string[], labels?: Labels): Field {
  return {
    name,
    type: FieldType.string,
    config: {},
    values,
    labels,
  };
}

function makeTime(name: string, values: number[], nanos?: number[]): Field {
  return {
    name,
    type: FieldType.time,
    config: {},
    values,
  };
}

function makeObject(name: string, values: Object[]): Field {
  return {
    name,
    type: FieldType.other,
    config: {},
    values,
  };
}

describe('parseLogsFrame should parse different logs-dataframe formats', () => {
  it('should parse a dataplane-complaint logs frame', () => {
    const time = makeTime('timestamp', [1687185711795, 1687185711995]);
    const body = makeString('body', ['line1', 'line2']);
    const severity = makeString('severity', ['info', 'debug']);
    const id = makeString('id', ['id1', 'id2']);
    const attributes = makeObject('attributes', [
      { counter: '38141', label: 'val2', level: 'warning' },
      { counter: '38143', label: 'val2', level: 'info' },
    ]);

    const result = parseLogsFrame({
      meta: {
        type: DataFrameType.LogLines,
      },
      fields: [id, body, attributes, severity, time],
      length: 2,
    });

    expect(result).not.toBeNull();
    if (result == null) {
      // to make typescript happy
      throw new Error('should never happen');
    }

    expect(result.timeField.values[0]).toBe(time.values[0]);
    expect(result.bodyField.values[0]).toBe(body.values[0]);
    expect(result.idField?.values[0]).toBe(id.values[0]);
    expect(result.timeNanosecondField).toBeUndefined();
    expect(result.severityField?.values[0]).toBe(severity.values[0]);
    expect(result.attributes).toStrictEqual([
      { counter: '38141', label: 'val2', level: 'warning' },
      { counter: '38143', label: 'val2', level: 'info' },
    ]);
  });

  it('should parse old Loki-style (grafana8.x) frames ( multi-frame, but here we only parse a single frame )', () => {
    const time = makeTime('ts', [1687185711795, 1687185711995]);
    const line = makeString('line', ['line1', 'line2'], { counter: '34543', lable: 'val3', level: 'info' });
    const id = makeString('id', ['id1', 'id2']);
    const ns = makeString('tsNs', ['1687185711795123456', '1687185711995987654']);

    const result = parseLogsFrame({
      fields: [time, line, ns, id],
      length: 2,
    });

    expect(result).not.toBeNull();
    if (result == null) {
      // to make typescript happy
      throw new Error('should never happen');
    }

    expect(result.timeField.values[0]).toBe(time.values[0]);
    expect(result.bodyField.values[0]).toBe(line.values[0]);
    expect(result.idField?.values[0]).toBe(id.values[0]);
    expect(result.timeNanosecondField?.values[0]).toBe(ns.values[0]);
    expect(result.severityField).toBeUndefined();
    expect(result.attributes).toStrictEqual([
      { counter: '34543', lable: 'val3', level: 'info' },
      { counter: '34543', lable: 'val3', level: 'info' },
    ]);
  });

  it('should parse a Loki-style frame (single-frame, labels-in-json)', () => {
    const time = makeTime('Time', [1687185711795, 1687185711995]);
    const line = makeString('Line', ['line1', 'line2']);
    const id = makeString('id', ['id1', 'id2']);
    const ns = makeString('tsNs', ['1687185711795123456', '1687185711995987654']);
    const labels = makeObject('labels', [
      { counter: '38141', label: 'val2', level: 'warning' },
      { counter: '38143', label: 'val2', level: 'info' },
    ]);

    const result = parseLogsFrame({
      meta: {
        custom: {
          frameType: 'LabeledTimeValues',
        },
      },
      fields: [labels, time, line, ns, id],
      length: 2,
    });

    expect(result).not.toBeNull();
    if (result == null) {
      // to make typescript happy
      throw new Error('should never happen');
    }

    expect(result.timeField.values[0]).toBe(time.values[0]);
    expect(result.bodyField.values[0]).toBe(line.values[0]);
    expect(result.idField?.values[0]).toBe(id.values[0]);
    expect(result.timeNanosecondField?.values[0]).toBe(ns.values[0]);
    expect(result.severityField).toBeUndefined();
    expect(result.attributes).toStrictEqual([
      { counter: '38141', label: 'val2', level: 'warning' },
      { counter: '38143', label: 'val2', level: 'info' },
    ]);
  });

  it('should parse elastic-style frame (has level-field, no labels parsed, extra fields ignored)', () => {
    const time = makeTime('Time', [1687185711795, 1687185711995]);
    const line = makeString('Line', ['line1', 'line2']);
    const source = makeObject('_source', [
      { counter: '38141', label: 'val2', level: 'warning' },
      { counter: '38143', label: 'val2', level: 'info' },
    ]);
    const host = makeString('hostname', ['h1', 'h2']);
    const level = makeString('level', ['info', 'error']);

    const result = parseLogsFrame({
      meta: {
        custom: {
          frameType: 'LabeledTimeValues',
        },
      },
      fields: [time, line, source, level, host],
      length: 2,
    });

    expect(result).not.toBeNull();
    if (result == null) {
      // to make typescript happy
      throw new Error('should never happen');
    }

    expect(result.timeField.values[0]).toBe(time.values[0]);
    expect(result.bodyField.values[0]).toBe(line.values[0]);
    expect(result.severityField?.values[0]).toBe(level.values[0]);
    expect(result.idField).toBeUndefined();
    expect(result.timeNanosecondField).toBeUndefined();
    expect(result.attributes).toBeUndefined();
  });
});
