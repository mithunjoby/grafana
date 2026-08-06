import { generatedAPI } from '@grafana/api-clients/rtkq/dashboard/v2beta2';

export const dashboardAPIv2beta2 = generatedAPI.enhanceEndpoints({});

export const { useGetNotebookQuery, useLazyGetNotebookQuery } = dashboardAPIv2beta2;
