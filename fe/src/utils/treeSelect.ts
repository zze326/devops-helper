import _ from 'lodash'
// 刷新树状选择框的选择状态
export const refershTreeSelectDataWithChecked = (data: ITreeSelectData[], checkedIDs: number[], allCheck: boolean = false): ITreeSelectData[] => {
    if (allCheck) {
        return data.map((item: ITreeSelectData) => {
            if (checkedIDs.includes(item.value ?? 0))
                checkedIDs.splice(checkedIDs.indexOf(item.value ?? 0), 1)
            return {
                ...item,
                checked: allCheck,
                children: item.children ? refershTreeSelectDataWithChecked(item.children, checkedIDs, allCheck) : undefined,
            }
        })
    } else {
        return data.map((item: ITreeSelectData) => {
            let checked = checkedIDs.includes(item.value ?? 0)
            return {
                ...item,
                checked: checked,
                children: item.children ? refershTreeSelectDataWithChecked(item.children, checkedIDs, checked) : undefined,
            }
        })
    }
}

export const genTreeSelectData = <T>(data: T[], path: T[] = [], props: { labelField: string, valueField: string, childrenField: string } = { labelField: "name", valueField: "id", childrenField: "children" }): ITreeSelectData[] => {
    return data.map((item: T) => {
        let labelWithPath = [...path, item].map(item => _.get(item, props.labelField)).join(' / ')
        let children = _.get(item, props.childrenField)
        return {
            label: _.get(item, props.labelField),
            value: _.get(item, props.valueField),
            children: children ? genTreeSelectData(children, [...path, item]) : undefined,
            labelWithPath,
        }
    })
}

export const genTreeSelectDataWithChecked = <T, K>(data: T[], checkedValues: K[], allCheck: boolean = false, path: T[] = [], props: { labelField: string, valueField: string, childrenField: string } = { labelField: "name", valueField: "id", childrenField: "children" }): ITreeSelectData[] => {
    if (allCheck) {
        return data.map((item: T) => {
            let labelWithPath = [...path, item].map(item => _.get(item, props.labelField)).join(' / ')
            let children = _.get(item, props.childrenField)
            return {
                label: _.get(item, props.labelField),
                value: _.get(item, props.valueField),
                checked: allCheck,
                children: children ? genTreeSelectDataWithChecked(children, checkedValues, allCheck, [...path, item]) : undefined,
                labelWithPath
            }
        })
    } else {
        return data.map((item: T) => {
            const labelWithPath = [...path, item].map(item => _.get(item, props.labelField)).join(' / ')
            let checked = checkedValues.includes(_.get(item, props.valueField))
            let value = _.get(item, props.valueField)
            let children = _.get(item, props.childrenField)
            return {
                label: _.get(item, props.labelField),
                value: value,
                checked: checkedValues.includes(value),
                children: children ? genTreeSelectDataWithChecked(children, checkedValues, checked, [...path, item]) : undefined,
                labelWithPath
            }
        })
    }
}