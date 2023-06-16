import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import { DataSourceInstanceSettings, DataSourcePluginMeta, PluginMetaInfo, PluginType } from '@grafana/data';
import { GetDataSourceListFilters } from '@grafana/runtime';

import { DataSourceModal } from './DataSourceModal';

const pluginMetaInfo: PluginMetaInfo = {
  author: { name: '' },
  description: '',
  screenshots: [],
  version: '',
  updated: '',
  links: [],
  logos: { small: '', large: '' },
};

function createPluginMeta(name: string, builtIn: boolean): DataSourcePluginMeta {
  return { builtIn, name, id: name, type: PluginType.datasource, baseUrl: '', info: pluginMetaInfo, module: '' };
}

function createDS(name: string, id: number, builtIn: boolean): DataSourceInstanceSettings {
  return {
    name: name,
    uid: name + 'uid',
    meta: createPluginMeta(name, builtIn),
    id,
    access: 'direct',
    jsonData: {},
    type: '',
    readOnly: true,
  };
}

const mockDS = createDS('mockDS', 1, false);
const xMockDS = createDS('xMockDS', 2, false);
const builtInMockDS = createDS('builtInMockDS', 3, true);

const mockDSList = [mockDS, xMockDS, builtInMockDS];

const setup = (onChange = () => {}, onDismiss = () => {}) => {
  const props = { onChange, onDismiss, current: mockDS.name };
  window.HTMLElement.prototype.scrollIntoView = function () {};
  return render(<DataSourceModal {...props}></DataSourceModal>);
};

jest.mock('@grafana/runtime', () => {
  const actual = jest.requireActual('@grafana/runtime');
  return {
    ...actual,
    getTemplateSrv: () => {
      return {
        getVariables: () => [],
      };
    },
  };
});

jest.mock('@grafana/runtime/src/services/dataSourceSrv', () => {
  return {
    getDataSourceSrv: () => ({
      getList: (filters: GetDataSourceListFilters) => {
        if (filters.filter) {
          return mockDSList.filter(filters.filter);
        }
        return mockDSList;
      },
      getInstanceSettings: () => mockDS,
    }),
  };
});

describe('DataSourceDropdown', () => {
  it('should render', () => {
    expect(() => setup()).not.toThrow();
  });

  describe('interactions', () => {
    const user = userEvent.setup();

    it('should be searchable', async () => {
      setup();
      const searchBox = await screen.findByRole('searchbox');
      expect(searchBox).toBeInTheDocument();
      await user.click(searchBox!);

      await user.keyboard(xMockDS.name); //Search for xMockDS

      expect(screen.queryByText(mockDS.name, { selector: 'span' })).toBeNull();
      expect(await screen.findByText(xMockDS.name, { selector: 'span' })).toBeInTheDocument();

      await user.keyboard('foobarbaz'); //Search for a DS that should not exist

      expect(await screen.findByText('No data sources found')).toBeInTheDocument();
    });

    it('should call the onChange handler with the correct datasource', async () => {
      const onChange = jest.fn();
      setup(onChange);
      await user.click(await screen.findByText(xMockDS.name, { selector: 'span' }));
      expect(onChange.mock.lastCall[0].name).toEqual(xMockDS.name);
    });
  });
});
