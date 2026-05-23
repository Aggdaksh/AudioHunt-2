import type { RecognitionStatus } from '../types';

interface RecognizeButtonProps {
  onPress: () => void;
  status: RecognitionStatus;
}

function copyForStatus(status: RecognitionStatus) {
  switch (status) {
    case 'recording':
      return {
        title: 'Listening...',
        hint: 'Press Stop when your clip is ready.',
      };
    case 'processing':
      return {
        title: 'Processing...',
        hint: 'Checking your clip against the library.',
      };
    case 'match':
      return {
        title: 'Start listening',
        hint: 'Tap to run another recording.',
      };
    case 'no_match':
      return {
        title: 'Start listening',
        hint: 'No match found. Try again.',
      };
    case 'error':
      return {
        title: 'Start listening',
        hint: 'Fix the recording and try again.',
      };
    default:
      return {
        title: 'Start listening',
        hint: 'Play music near the mic for 8 to 10 seconds, then press Stop.',
      };
  }
}

export default function RecognizeButton({ onPress, status }: RecognizeButtonProps) {
  const isBusy = status === 'recording' || status === 'processing';
  const copy = copyForStatus(status);

  return (
    <button
      aria-label={status === 'recording' ? 'Listening' : 'Start recording'}
      className={`recognize-orb recognize-orb--${status}`}
      disabled={isBusy}
      onClick={onPress}
      type="button"
    >
      <span className="recognize-orb__halo" aria-hidden="true" />
      <span className="recognize-orb__halo recognize-orb__halo--delay" aria-hidden="true" />
      <span className="recognize-orb__core">
        <span className="recognize-orb__glyph" aria-hidden="true">
          <>
            <span />
            <span />
            <span />
          </>
        </span>
        <span className="recognize-orb__title">{copy.title}</span>
        <span className="recognize-orb__hint">{copy.hint}</span>
      </span>
      {status === 'processing' ? <span className="recognize-orb__spinner" aria-hidden="true" /> : null}
    </button>
  );
}
