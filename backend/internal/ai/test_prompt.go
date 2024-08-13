package ai

const testPrompt = `
<tool id="polaroid">
<file name="tool.ts">
import { BaseBoxShapeTool, TLClickEvent } from 'tldraw'

export class PolaroidShapeTool extends BaseBoxShapeTool {
        static override id = 'polaroid'
        static override initial = 'idle'
        override shapeType = 'polaroid'

        override onDoubleClick: TLClickEvent = (info) => {
                const shape = this.editor.getShapeByElement(info.target)
                if (shape && shape.type === 'polaroid') {
                        // Open a file picker dialog when double-clicked
                        const input = document.createElement('input')
                        input.type = 'file'
                        input.accept = 'image/*'
                        input.onchange = (e: Event) => {
                                const file = (e.target as HTMLInputElement).files?.[0]
                                if (file) {
                                        const reader = new FileReader()
                                        reader.onload = (e) => {
                                                const imageUrl = e.target?.result as string
                                                this.editor.updateShape({
                                                        id: shape.id,
                                                        type: 'polaroid',
                                                        props: { ...shape.props, imageUrl },
                                                })
                                        }
                                        reader.readAsDataURL(file)
                                }
                        }
                        input.click()
                }
        }
}
</file>

<file name="util.tsx">
import { useState, useEffect } from 'react'
import {
        HTMLContainer,
        Rectangle2d,
        ShapeUtil,
        TLOnResizeHandler,
        getDefaultColorTheme,
        resizeBox,
} from 'tldraw'

interface IPolaroidShape {
        type: 'polaroid'
        props: {
                w: number
                h: number
                title: string
                imageUrl: string
        }
}

export class PolaroidShapeUtil extends ShapeUtil<IPolaroidShape> {
        static override type = 'polaroid' as const

        getDefaultProps(): IPolaroidShape['props'] {
                return {
                        w: 200,
                        h: 240,
                        title: 'New Polaroid',
                        imageUrl: '',
                }
        }

        getGeometry(shape: IPolaroidShape) {
                return new Rectangle2d({
                        width: shape.props.w,
                        height: shape.props.h,
                        isFilled: true,
                })
        }

        component(shape: IPolaroidShape) {
                const bounds = this.editor.getShapeGeometry(shape).bounds
                const theme = getDefaultColorTheme({ isDarkMode: this.editor.user.getIsDarkMode() })

                // eslint-disable-next-line react-hooks/rules-of-hooks
                const [title, setTitle] = useState(shape.props.title)

                // eslint-disable-next-line react-hooks/rules-of-hooks
                useEffect(() => {
                        setTitle(shape.props.title)
                }, [shape.props.title])

                const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
                        const newTitle = e.target.value
                        setTitle(newTitle)
                        this.editor.updateShape({
                                id: shape.id,
                                type: 'polaroid',
                                props: { ...shape.props, title: newTitle },
                        })
                }

                return (
                        <HTMLContainer
                                id={shape.id}
                                style={{
                                        width: '100%',
                                        height: '100%',
                                        backgroundColor: 'white',
                                        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
                                        display: 'flex',
                                        flexDirection: 'column',
                                        padding: '10px',
                                        boxSizing: 'border-box',
                                }}
                        >
                                <div
                                        style={{
                                                flex: 1,
                                                backgroundColor: '#f0f0f0',
                                                marginBottom: '10px',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                                overflow: 'hidden',
                                        }}
                                >
                                        {shape.props.imageUrl ? (
                                                <img
                                                        src={shape.props.imageUrl}
                                                        alt="Polaroid"
                                                        style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }}
                                                />
                                        ) : (
                                                <span>Double-click to add image</span>
                                        )}
                                </div>
                                <input
                                        type="text"
                                        value={title}
                                        onChange={handleTitleChange}
                                        style={{
                                                width: '100%',
                                                border: 'none',
                                                background: 'transparent',
                                                textAlign: 'center',
                                                fontSize: '14px',
                                                fontFamily: 'Arial, sans-serif',
                                        }}
                                        onPointerDown={(e) => e.stopPropagation()}
                                />
                        </HTMLContainer>
                )
        }

        indicator(shape: IPolaroidShape) {
                return <rect width={shape.props.w} height={shape.props.h} />
        }

        override onResize: TLOnResizeHandler<IPolaroidShape> = (shape, info) => {
                return resizeBox(shape, info)
        }
}
</file>

<file name="icon.svg">
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
  <rect x="5" y="5" width="14" height="10"/>
  <line x1="5" y1="18" x2="19" y2="18"/>
</svg>
</file>
</tool>
`
