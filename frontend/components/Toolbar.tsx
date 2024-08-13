/* eslint-disable @next/next/no-img-element */
import { Hand, MousePointer2, PencilIcon } from "lucide-react";
import { Card } from "@/components/ui/card";
import { ToggleGroupItem, ToggleGroup } from "@/components/ui/toggle-group";
import { TLStateNodeConstructor, useEditor, useValue } from "tldraw";
import { getToolIconId, toolIcons } from "@/components/tldraw-custom-tools";

type Props = {
  customTools: TLStateNodeConstructor[];
};

export const Toolbar = ({ customTools }: Props) => {
  const editor = useEditor();

  const selectedTool = useValue("tool", () => editor.getCurrentToolId(), [
    editor,
  ]);

  return (
    <Card className="p-1">
      <ToggleGroup
        type="single"
        size="sm"
        onValueChange={(v) => {
          editor.setCurrentTool(v);
        }}
        value={selectedTool}
      >
        <ToggleGroupItem value="select" aria-label="Select" title="Select">
          <MousePointer2 className="h-4 w-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="draw" aria-label="draw" title="Draw">
          <PencilIcon className="h-4 w-4" />
        </ToggleGroupItem>
        <ToggleGroupItem value="hand" aria-label="hand" title="Pan">
          <Hand className="h-4 w-4" />
        </ToggleGroupItem>

        {customTools.map((tool) => (
          <ToggleGroupItem
            key={tool.id}
            value={tool.id}
            aria-label={tool.id}
            title={tool.name}
          >
            <img
              src={toolIcons[getToolIconId(tool.id)]}
              className="w-4 h-4"
              alt={getToolIconId(tool.id)}
            />
          </ToggleGroupItem>
        ))}
      </ToggleGroup>
    </Card>
  );
};
