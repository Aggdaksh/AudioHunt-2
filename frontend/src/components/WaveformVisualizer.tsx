interface WaveformVisualizerProps {
  frequencyData: Uint8Array;
  visible: boolean;
}

function buildBars(frequencyData: Uint8Array, count: number) {
  const bars = Array.from({ length: count }, () => 0.08);

  if (!frequencyData.length) {
    return bars;
  }

  const binsPerBar = Math.max(1, Math.floor(frequencyData.length / count));

  for (let index = 0; index < count; index += 1) {
    const start = index * binsPerBar;
    const end = Math.min(frequencyData.length, start + binsPerBar);
    let total = 0;

    for (let cursor = start; cursor < end; cursor += 1) {
      total += frequencyData[cursor];
    }

    const average = total / Math.max(1, end - start);
    bars[index] = Math.max(0.08, average / 255);
  }

  return bars;
}

export default function WaveformVisualizer({ frequencyData, visible }: WaveformVisualizerProps) {
  const bars = buildBars(frequencyData, 40);

  return (
    <div className={`waveform ${visible ? 'waveform--visible' : ''}`} aria-hidden="true">
      {bars.map((bar, index) => (
        <span
          className="waveform__bar"
          key={`${index}-${Math.round(bar * 1000)}`}
          style={{ '--bar-scale': `${bar}` } as React.CSSProperties}
        />
      ))}
    </div>
  );
}
