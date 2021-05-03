import { matchPath } from 'react-router-dom';
import qs, { ParsedQs } from 'qs';
import { Location, LocationState } from 'history';

import useCases from 'constants/useCaseTypes';
import { searchParams, sortParams, pagingParams } from 'constants/searchParams';
import WorkflowEntity from './WorkflowEntity';
import { WorkflowState } from './WorkflowState';
import {
    workflowPaths,
    urlEntityListTypes,
    urlEntityTypes,
    clustersPathWithParam,
    riskPath,
    violationsPath,
    policiesPath,
    networkPath,
    userRolePath,
} from '../routePaths';

function getTypeKeyFromParamValue(value: string, listOnly = false): string | null {
    const listMatch = Object.entries(urlEntityListTypes).find((entry) => entry[1] === value);
    const entityMatch = Object.entries(urlEntityTypes).find((entry) => entry[1] === value);
    const match = listOnly ? listMatch : listMatch || entityMatch;
    return match ? match[0] : null;
}

function getEntityFromURLParam(type: string, id?: string): WorkflowEntity {
    return new WorkflowEntity(getTypeKeyFromParamValue(type), id);
}

function paramsToStateStack(params): WorkflowEntity[] {
    const { pageEntityListType, pageEntityType, pageEntityId, entityId1, entityId2 } = params;
    const { entityType1: urlEntityType1, entityType2: urlEntityType2 } = params;
    const entityListType1 = getTypeKeyFromParamValue(urlEntityType1, true);
    const entityListType2 = getTypeKeyFromParamValue(urlEntityType2, true);
    const entityType1 = getTypeKeyFromParamValue(urlEntityType1);
    const entityType2 = getTypeKeyFromParamValue(urlEntityType2);
    const stateArray: WorkflowEntity[] = [];
    if (!pageEntityListType && !pageEntityType) {
        return stateArray;
    }

    // List
    if (pageEntityListType) {
        stateArray.push(getEntityFromURLParam(pageEntityListType));

        if (entityId1) {
            stateArray.push(getEntityFromURLParam(pageEntityListType, entityId1));
        }
    } else {
        stateArray.push(getEntityFromURLParam(pageEntityType, pageEntityId));
        if (entityListType1) {
            stateArray.push(new WorkflowEntity(entityListType1));
        }
        if (entityType1 && entityId1) {
            stateArray.push(new WorkflowEntity(entityType1, entityId1));
        }
    }

    if (entityListType2) {
        stateArray.push(new WorkflowEntity(entityListType2));
    }
    if (entityType2 && entityId2) {
        stateArray.push(new WorkflowEntity(entityType2, entityId2));
    }

    return stateArray;
}

function formatSort(sort?: ParsedQs | ParsedQs[]): Record<string, unknown>[] | null {
    if (!sort) {
        return null;
    }

    let sorts: ParsedQs[];
    if (!Array.isArray(sort)) {
        sorts = [sort];
    } else {
        sorts = [...sort];
    }

    return sorts.map(({ id, desc }) => {
        return {
            id,
            desc: JSON.parse(desc as string),
        } as Record<string, unknown>;
    });
}

// Convert URL to workflow state and search objects
// note: this will read strictly from 'location' as 'match' is relative to the closest Route component
function parseURL(location: Location<LocationState>): WorkflowState {
    if (!location) {
        // TODO: be more specific, it could be an exception instead of a dummy object
        return new WorkflowState();
    }

    const { pathname, search } = location;
    const listParams = matchPath(pathname, {
        path: workflowPaths.LIST,
    });
    const entityParams = matchPath(pathname, {
        path: workflowPaths.ENTITY,
    });
    const dashboardParams = matchPath(pathname, {
        path: workflowPaths.DASHBOARD,
        exact: true,
    });

    // check for legacy-Workflow sections
    const matchedRiskParams = matchPath(pathname, {
        path: riskPath,
        exact: true,
    });
    const matchedNetworkParams = matchPath(pathname, {
        path: networkPath,
        exact: true,
    });
    const matchedClustersParams = matchPath(pathname, {
        path: clustersPathWithParam,
        exact: true,
    });
    const matchedViolationsParams = matchPath(pathname, {
        path: violationsPath,
        exact: true,
    });
    const matchedPoliciesParams = matchPath(pathname, {
        path: policiesPath,
        exact: true,
    });
    let legacyParams = { params: {} };
    if (matchedNetworkParams) {
        legacyParams = {
            params: {
                ...matchedNetworkParams,
                context: useCases.NETWORK,
            },
        };
    }
    if (matchedRiskParams) {
        legacyParams = {
            params: {
                ...matchedRiskParams,
                context: useCases.RISK,
            },
        };
    }
    if (matchedClustersParams) {
        legacyParams = {
            params: {
                ...matchedClustersParams,
                context: useCases.CLUSTERS,
            },
        };
    }
    if (matchedViolationsParams) {
        legacyParams = {
            params: {
                ...matchedViolationsParams,
                context: useCases.VIOLATIONS,
            },
        };
    }
    if (matchedPoliciesParams) {
        legacyParams = {
            params: {
                ...matchedPoliciesParams,
                context: useCases.POLICIES,
            },
        };
    }

    const matchedUserRoleParams = matchPath(pathname, {
        path: userRolePath,
        exact: true,
    });
    if (matchedUserRoleParams) {
        legacyParams = {
            params: {
                ...matchedUserRoleParams,
                context: useCases.USER,
            },
        };
    }

    const { params } = entityParams || listParams || dashboardParams || legacyParams;
    const queryStr = search ? qs.parse(search, { ignoreQueryPrefix: true }) : {};

    const stateStackFromURLParams = paramsToStateStack(params) || [];

    const {
        [searchParams.page]: pageSearch,
        [searchParams.sidePanel]: sidePanelSearch,
        [sortParams.page]: pageSort,
        [sortParams.sidePanel]: sidePanelSort,
        [pagingParams.page]: pagePaging,
        [pagingParams.sidePanel]: sidePanelPaging,
    } = queryStr;

    const queryWorkflowState = queryStr.workflowState || [];
    const stateStackFromQueryString = !Array.isArray(queryWorkflowState)
        ? [queryWorkflowState as ParsedQs]
        : (queryWorkflowState as ParsedQs[]);
    const stateStack = stateStackFromQueryString.map(({ t, i }) => new WorkflowEntity(t, i));

    const workflowState = new WorkflowState(
        params.context,
        [...stateStackFromURLParams, ...stateStack],
        {
            [searchParams.page]: pageSearch || null,
            [searchParams.sidePanel]: sidePanelSearch || null,
        },
        {
            [sortParams.page]: formatSort(pageSort as ParsedQs | ParsedQs[]),
            [sortParams.sidePanel]: formatSort(sidePanelSort as ParsedQs | ParsedQs[]),
        },
        {
            [pagingParams.page]: parseInt((pagePaging as string) ?? '0', 10),
            [pagingParams.sidePanel]: parseInt((sidePanelPaging as string) ?? '0', 10),
        }
    );

    return workflowState;
}

export default parseURL;
