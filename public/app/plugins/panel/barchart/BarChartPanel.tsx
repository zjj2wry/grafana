import React, { PureComponent } from 'react';
import { PanelProps, getFieldDisplayName } from '@grafana/data';
import { BarChartOptions } from './types';

export class BarChartPanel extends PureComponent<PanelProps<BarChartOptions>> {
  render() {
    const { options, data } = this.props;

    if (data.series) {
      for (const frame of data.series) {
        for (const field of frame.fields) {
          const name = getFieldDisplayName(field, frame);
          console.log(name);
        }
      }
    }

    return (
      <div>
        Hello Chart!!!
        <pre>{JSON.stringify(options)}</pre>
      </div>
    );
  }
}
