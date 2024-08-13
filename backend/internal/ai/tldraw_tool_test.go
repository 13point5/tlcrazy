package ai

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	exampleToolId = "sticker"

	exampleToolFile = `
import { StateNode } from "tldraw";

const OFFSET = 12;
class StickerTool extends StateNode {
  static override id = 'sticker';

  override onEnter = () => {
    this.editor.setCursor({ type: "cross", rotation: 0 });
  };

  override onPointerDown = () => {
    const { currentPagePoint } = this.editor.inputs;
    this.editor.createShape({
      type: "text",
      x: currentPagePoint.x - OFFSET,
      y: currentPagePoint.y - OFFSET,
      props: { text: "❤️" },
    });
  };
}

export default StickerTool;
`

	exampleToolIcon = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <path
    d="M50 10 L61 38 90 38 67 56 76 84 50 68 24 84 33 56 10 38 39 38Z"
    stroke="black"
    stroke-width="2"
    fill="none"
  />
</svg>
`
)

func TestParseTldrawToolRawOutput(t *testing.T) {

	testcase := TldrawToolOutput{
		Id:   exampleToolId,
		Tool: exampleToolFile,
		Icon: exampleToolIcon,
	}

	toolXML := fmt.Sprintf(`
<tool id="%s">
<file name="tool.ts">%s</file>

<file name="icon.svg">%s</file>
</tool>
`, testcase.Id, testcase.Tool, testcase.Icon)

	t.Run("Check if ID, Tool, and Icons are extracted correctly", func(t *testing.T) {
		out, err := parseTldrawToolXML(toolXML)
		if err != nil {
			t.Fatal("Got an error but didn't expect one", err)
		}

		if !reflect.DeepEqual(testcase, out) {
			t.Errorf("Expected %q\nbut got %q", testcase, out)
		}
	})
}
