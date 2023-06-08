import { DataFrame, FieldCache, FieldType, FieldWithIndex, Labels, DataFrameType } from '@grafana/data';

import { parseLegacyLogsFrame } from './legacyLogsFrame';

// these are like Labels, but their values can be
// arbitrary structures, not just strings
export type Attributes = Record<string, any>;

// the attributes-access is a little awkward, but it's necessary
// because there are multiple,very different dataframe-represenations.
export type LogsFrame = {
  timeField: FieldWithIndex;
  bodyField: FieldWithIndex;
  timeNanosecondField?: FieldWithIndex;
  severityField?: FieldWithIndex;
  idField?: FieldWithIndex;
  attributes?: Attributes[];
};

function getField(cache: FieldCache, name: string, fieldType: FieldType): FieldWithIndex | undefined {
  const field = cache.getFieldByName(name);
  if (field == null) {
    return undefined;
  }

  return field.type === fieldType ? field : undefined;
}

const DATAPLANE_TIMESTAMP_NAME = 'timestamp';
const DATAPLANE_BODY_NAME = 'body';
const DATAPLANE_SEVERITY_NAME = 'severity';
const DATAPLANE_ID_NAME = 'id';
const DATAPLANE_ATTRIBUTES_NAME = 'attributes';

function parseDataplaneLogsFrame(frame: DataFrame): LogsFrame | null {
  const cache = new FieldCache(frame);

  const timestampField = getField(cache, DATAPLANE_TIMESTAMP_NAME, FieldType.time);
  const bodyField = getField(cache, DATAPLANE_BODY_NAME, FieldType.string);

  // these two are mandatory
  if (timestampField == null || bodyField == null) {
    return null;
  }

  const severityField = getField(cache, DATAPLANE_SEVERITY_NAME, FieldType.string);
  const idField = getField(cache, DATAPLANE_ID_NAME, FieldType.string);
  const attributesField = getField(cache, DATAPLANE_ATTRIBUTES_NAME, FieldType.other);

  const attributes = attributesField == null ? undefined : attributesField.values;

  return {
    timeField: timestampField,
    bodyField,
    severityField,
    idField,
    attributes,
  };
  return null;
}

export function parseLogsFrame(frame: DataFrame): LogsFrame | null {
  if (frame.meta?.type === DataFrameType.LogLines) {
    return parseDataplaneLogsFrame(frame);
  } else {
    return parseLegacyLogsFrame(frame);
  }
}
