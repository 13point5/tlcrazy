
import React, { useState } from 'react'
import {
	HTMLContainer,
	Rectangle2d,
	ShapeUtil,
	TLOnResizeHandler,
	resizeBox,
} from 'tldraw'

interface YouTubePlayerShape {
	type: 'youtube-player'
	props: {
		w: number
		h: number
		url: string
	}
}

export default class YouTubePlayerUtil extends ShapeUtil<YouTubePlayerShape> {
	static override type = 'youtube-player' as const

	override isAspectRatioLocked = (_shape: YouTubePlayerShape) => false
	override canResize = (_shape: YouTubePlayerShape) => true

	getDefaultProps(): YouTubePlayerShape['props'] {
		return {
			w: 560,
			h: 315,
			url: '',
		}
	}

	getGeometry(shape: YouTubePlayerShape) {
		return new Rectangle2d({
			width: shape.props.w,
			height: shape.props.h,
			isFilled: true,
		})
	}

	component(shape: YouTubePlayerShape) {
		const [url, setUrl] = useState(shape.props.url)

		const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
			const newUrl = e.target.value
			setUrl(newUrl)
			this.editor.updateShape<YouTubePlayerShape>({
				id: shape.id,
				type: 'youtube-player',
				props: { ...shape.props, url: newUrl },
			})
		}

		const embedUrl = url.replace('watch?v=', 'embed/')

		return (
			<HTMLContainer
				id={shape.id}
				style={{
					pointerEvents: 'all',
				}}
			>
				<div className="flex flex-col w-full h-full bg-white">
					<input
						type="text"
						value={url}
						onChange={handleUrlChange}
						placeholder="Enter YouTube URL"
						className="p-2 mb-2 border border-gray-300 rounded"
						onPointerDown={(e) => e.stopPropagation()}
					/>
					{url && (
						<iframe
							width="100%"
							height="100%"
							src={embedUrl}
							frameBorder="0"
							allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
							allowFullScreen
						></iframe>
					)}
				</div>
			</HTMLContainer>
		)
	}

	indicator(shape: YouTubePlayerShape) {
		return <rect width={shape.props.w} height={shape.props.h} />
	}

	override onResize: TLOnResizeHandler<YouTubePlayerShape> = (shape, info) => {
		return resizeBox(shape, info)
	}
}
