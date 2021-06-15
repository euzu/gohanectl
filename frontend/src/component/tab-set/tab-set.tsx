import React, {forwardRef, useImperativeHandle, useMemo, useState} from 'react';
import './tab-set.scss';
import {Tab, Tabs} from "@material-ui/core";
import useTranslator from "../../hook/use-translator";

interface TabDescription {
    label: string;
    icon?: React.ReactElement;
    children?: React.ReactNode;
}

interface TabSetProps {
    tabs: TabDescription[];
    index?: number;
    rolling?: boolean; // swipe rolling
    onChange?: (index: number) => void;
}

interface TabPanelProps {
    children?: React.ReactNode;
    index: any;
    value: any;
}

function TabPanel(props: TabPanelProps) {
    const {children, value, index, ...other} = props;
    return (
        <div
            className={'tabset-panel ' + (value === index ? 'tab-active' : 'tab-hidden')}
            role="tabpanel"
            hidden={value !== index}
            //aria-labelledby={`tabPanel-${index}`}
            {...other}>
            {children}
        </div>
    );
}

export interface ITabSet {
    setActiveTab(idx: number): void;

    setNextTab(): void;

    setPrevTab(): void;
}

const TabSet = forwardRef<ITabSet, TabSetProps>((props, ref) => {
    const {tabs, index, rolling, onChange} = props;
    const maxTabIndex = tabs.length - 1;
    const translate = useTranslator();
    const [activeTab, setActiveTab] = useState(index || 0);
    const renderTabs = useMemo(() => tabs.map((tab: TabDescription, idx: number) =>
        <Tab icon={tab.icon} key={'tab' + idx} label={translate(tab.label)}/>), [tabs, translate]);
    const renderTabPanels = useMemo(() => tabs.filter((tab: TabDescription) => tab.children).map((tab: TabDescription, idx: number) =>
        <TabPanel key={'tab-panel' + idx} value={activeTab} index={idx}>{tab.children}</TabPanel>), [tabs, activeTab])
    const handleChange = useMemo(() => (event: React.ChangeEvent<{}>, newValue: number) => {
        setActiveTab(value => newValue);
        onChange && onChange(newValue);
    }, [onChange]);
    useImperativeHandle(ref, () => ({
        setActiveTab: (idx: number) => setActiveTab(value => idx),
        setNextTab: () => {
            let newIndex = activeTab;
            if (rolling && activeTab === maxTabIndex) {
                newIndex = 0;
            } else {
                newIndex = Math.min(activeTab + 1, maxTabIndex);
            }
            setActiveTab(() => newIndex)
            onChange && newIndex !== activeTab && onChange(newIndex);
        },
        setPrevTab: () => {
            let newIndex = activeTab;
            if (rolling && activeTab === 0) {
                newIndex = maxTabIndex;
            } else {
                newIndex = Math.max(activeTab - 1, 0)
            }
            setActiveTab(() => newIndex);
            onChange && newIndex !== activeTab && onChange(newIndex);
        },
    }));

    return <div className={'tabset-container'}>
        <Tabs value={activeTab} onChange={handleChange}>
            {renderTabs}
        </Tabs>
        {renderTabPanels}
    </div>
})

export default TabSet;
