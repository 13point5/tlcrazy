import {
  TLUiAssetUrlOverrides,
  TLStateNodeConstructor,
  TLAnyShapeUtilConstructor,
} from "tldraw";
import tools from "@/components/tldraw-custom-tools/tools.json";

const toolIds = tools.ids;

const importTool = async (toolId: string): Promise<TLStateNodeConstructor> => {
  const toolModule = await import(`./${toolId}/tool`);
  return toolModule.default;
};

export const loadTools = async (): Promise<TLStateNodeConstructor[]> => {
  return await Promise.all(toolIds.map(importTool));
};

const importToolUtil = async (
  toolId: string
): Promise<TLAnyShapeUtilConstructor> => {
  const utilModule = await import(`./${toolId}/util`);
  return utilModule.default;
};

export const loadUtils = async (): Promise<TLAnyShapeUtilConstructor[]> => {
  return await Promise.all(toolIds.map(importToolUtil));
};

export const getToolIconId = (toolId: string) => `${toolId}-tool-icon`;

export const toolIcons: Exclude<TLUiAssetUrlOverrides["icons"], undefined> =
  toolIds.reduce((acc, toolName) => {
    acc[getToolIconId(toolName)] = `/custom-tool-icons/${toolName}.svg`;
    return acc;
  }, {} as Record<string, string>);
