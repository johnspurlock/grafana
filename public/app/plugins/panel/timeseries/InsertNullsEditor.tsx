import React from 'react';

import { FieldOverrideEditorProps, rangeUtil, SelectableValue } from '@grafana/data';
import { HorizontalGroup, Input, RadioButtonGroup } from '@grafana/ui';

const DISCONNECT_OPTIONS: Array<SelectableValue<boolean | number>> = [
  {
    label: 'Never',
    value: false,
  },
  {
    label: 'Threshold',
    value: 3600000, // 1h
  },
];

type Props = FieldOverrideEditorProps<boolean | number, unknown>;

export const InsertNullsEditor = ({ value, onChange }: Props) => {
  const isThreshold = typeof value === 'number';
  const formattedTime = isThreshold ? rangeUtil.secondsToHms(value / 1000) : undefined;
  DISCONNECT_OPTIONS[1].value = isThreshold ? value : 3600000; // 1h

  const checkAndUpdate = (txt: string) => {
    let val: boolean | number = false;
    if (txt) {
      try {
        val = rangeUtil.intervalToMs(txt);
      } catch (err) {
        console.warn('ERROR', err);
      }
    }
    onChange(val);
  };

  const handleEnterKey = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== 'Enter') {
      return;
    }
    checkAndUpdate(e.currentTarget.value);
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    checkAndUpdate(e.currentTarget.value);
  };

  return (
    <HorizontalGroup>
      <RadioButtonGroup value={value} options={DISCONNECT_OPTIONS} onChange={onChange} />
      {isThreshold && (
        <Input
          autoFocus={false}
          placeholder="never"
          width={10}
          defaultValue={formattedTime}
          onKeyDown={handleEnterKey}
          onBlur={handleBlur}
          prefix={<div>&gt;</div>}
          spellCheck={false}
        />
      )}
    </HorizontalGroup>
  );
};
