import { PanelPlugin } from '@grafana/data';
import { BarChartPanel } from './BarChartPanel';
import { BarChartOptions } from './types';

export const plugin = new PanelPlugin<BarChartOptions>(BarChartPanel)
  .useFieldConfig() // Want data configs
  .setPanelOptions(builder => {
    builder
      .addNumberInput({
        path: 'minSpace',
        name: 'Min Space',
        defaultValue: 0,
      })
      .addNumberInput({
        path: 'maxSpace',
        name: 'Max Space',
        defaultValue: 1,
      });
  });
