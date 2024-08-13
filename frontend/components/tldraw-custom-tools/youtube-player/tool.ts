
import { BaseBoxShapeTool } from 'tldraw'

export default class YouTubePlayerTool extends BaseBoxShapeTool {
	static override id = 'youtube-player'
	static override initial = 'idle'
	override shapeType = 'youtube-player'
}
