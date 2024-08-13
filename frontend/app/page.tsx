"use client";

import {
  toolIcons,
  loadTools,
  loadUtils,
} from "@/components/tldraw-custom-tools";
import { Toolbar } from "@/components/Toolbar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import axios from "axios";
import { Loader2Icon, PencilIcon } from "lucide-react";
import { useEffect, useState } from "react";
import {
  DefaultToolbar,
  DefaultToolbarContent,
  TLAnyShapeUtilConstructor,
  TLComponents,
  TLStateNodeConstructor,
  Tldraw,
  TldrawUiMenuItem,
  useEditor,
} from "tldraw";
import "tldraw/tldraw.css";
import toolsData from "@/components/tldraw-custom-tools/tools.json";

const numOfTools = toolsData.ids.length;

export default function Home() {
  const [loadedTools, setLoadedTools] = useState<TLStateNodeConstructor[]>([]);
  const [loadedUtils, setLoadedUtils] = useState<TLAnyShapeUtilConstructor[]>(
    []
  );
  const [loadedComponents, setLoadedComponents] = useState<TLComponents | null>(
    null
  );
  const [query, setQuery] = useState("");
  const [creatingTool, setCreatingTool] = useState(false);

  const loadAllResources = async () => {
    const [tools, utils] = await Promise.all([loadTools(), loadUtils()]);

    setLoadedTools(tools);
    setLoadedComponents(loadComponents(tools));
    setLoadedUtils(utils);
  };

  useEffect(() => {
    loadAllResources();
  }, [toolsData.ids]);

  const handleCreateTool = async () => {
    setCreatingTool(true);

    try {
      const res = await axios.post("http://localhost:8080/tldraw-tool", {
        query,
      });
      console.log("res.data", res.data);
    } catch (error) {
      console.error(error);
    } finally {
      setCreatingTool(false);
    }
  };

  console.log({
    loadedTools,
    loadedUtils,
    loadedComponents,
    numOfTools,
    toolsData,
  });

  if (
    loadedTools.length !== numOfTools ||
    loadedUtils.length !== numOfTools ||
    loadedComponents === null
  ) {
    return <div>Loading...</div>;
  }

  return (
    <div style={{ position: "fixed", inset: 0 }}>
      <Tldraw
        tools={loadedTools}
        shapeUtils={loadedUtils}
        assetUrls={{
          icons: toolIcons,
        }}
        components={{
          Toolbar: null,
        }}
      >
        <div className="absolute top-2 left-1/2 transform -translate-x-1/2 z-[1000]">
          <div className="flex flex-col items-center gap-4">
            <h1 className="text-3xl font-medium">Make a Tool</h1>

            <div className="flex items-start gap-2">
              <Textarea
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
              <Button
                onClick={handleCreateTool}
                size="icon"
                className="bg-blue-600 hover:bg-blue-700"
              >
                {creatingTool ? (
                  <Loader2Icon className="animate-spin" size={16} />
                ) : (
                  <PencilIcon size={16} />
                )}
              </Button>
            </div>
          </div>
        </div>

        <div
          className="absolute bottom-1 left-1/2 transform -translate-x-1/2"
          style={{ zIndex: 1000 }}
        >
          <Toolbar customTools={loadedTools} />
        </div>
      </Tldraw>
    </div>
  );
}

const loadComponents = (tools: TLStateNodeConstructor[]): TLComponents => {
  return {
    Toolbar: (props) => {
      const editor = useEditor();
      const currentToolId = editor.getCurrentToolId();
      console.log("currentToolId", currentToolId);

      return (
        <DefaultToolbar {...props}>
          {tools.map((tool) => (
            <TldrawUiMenuItem
              key={tool.id}
              {...tool}
              onSelect={() => {
                editor.setCurrentTool(tool.id);
              }}
              isSelected={tool.id === currentToolId}
            />
          ))}

          <DefaultToolbarContent />
        </DefaultToolbar>
      );
    },
  };
};
