import React from 'react';
import {
  AbsoluteFill, Composition, Html5Audio, OffthreadVideo, Sequence,
  interpolate, registerRoot, staticFile, useCurrentFrame, useVideoConfig,
} from 'remotion';
import defaults from '../props.json';

const metadata = (props) => ({
  durationInFrames: props.shots.reduce((total, shot) => total + shot.durationFrames, 0),
  fps: props.fps,
  width: props.width,
  height: props.height,
});

const Shot = ({shot, title, index, total}) => {
  const frame = useCurrentFrame();
  const {width, height} = useVideoConfig();
  const vertical = height > width;
  const unit = Math.min(width, height) / 720;
  const entrance = Math.min(12, Math.max(1, shot.durationFrames - 1));
  const opacity = interpolate(frame, [0, entrance], [0, 1], {extrapolateRight: 'clamp'});
  const offset = interpolate(frame, [0, entrance], [18, 0], {extrapolateRight: 'clamp'});
  return (
    <AbsoluteFill style={{background: '#102234', color: '#f4f7fa', fontFamily: 'Arial, sans-serif', padding: 64 * unit}}>
      <div style={{fontSize: 24 * unit, color: '#80ded0', letterSpacing: 1}}>{title}</div>
      <div style={{flex: 1, display: 'flex', flexDirection: vertical ? 'column' : 'row', alignItems: 'center', gap: 32 * unit, opacity, transform: `translateY(${offset * unit}px)`}}>
        <div style={{flex: 1, width: '100%'}}>
          <div style={{fontSize: 18 * unit, color: '#80ded0', marginBottom: 20 * unit}}>0{index + 1} / 0{total}</div>
          <h1 style={{fontSize: (vertical ? 62 : 68) * unit, lineHeight: 1.08, margin: 0, overflowWrap: 'break-word'}}>{shot.headline}</h1>
          <p style={{fontSize: 26 * unit, lineHeight: 1.4, color: '#bacbd9', marginBottom: 0}}>{shot.body}</p>
        </div>
        {shot.media ? <div style={{width: vertical ? '100%' : '42%', height: vertical ? '38%' : '70%', overflow: 'hidden', borderRadius: 16 * unit}}>
          <OffthreadVideo src={staticFile(shot.media)} muted style={{width: '100%', height: '100%', objectFit: 'contain'}} />
        </div> : null}
      </div>
      <div style={{fontSize: 23 * unit, lineHeight: 1.4, minHeight: 64 * unit, paddingTop: 16 * unit, borderTop: `${2 * unit}px solid #365367`}}>{shot.caption || ''}</div>
    </AbsoluteFill>
  );
};

const ProductStory = (props) => {
  let cursor = 0;
  return <AbsoluteFill>
    {props.shots.map((shot, index) => {
      const from = cursor;
      cursor += shot.durationFrames;
      return <Sequence key={shot.id} from={from} durationInFrames={shot.durationFrames}>
        <Shot shot={shot} title={props.title} index={index} total={props.shots.length} />
      </Sequence>;
    })}
    {props.audio ? <Html5Audio src={staticFile(props.audio)} /> : null}
  </AbsoluteFill>;
};

const Root = () => <Composition id="ProductStory" component={ProductStory} defaultProps={defaults}
  {...metadata(defaults)} calculateMetadata={({props}) => metadata(props)} />;

registerRoot(Root);
