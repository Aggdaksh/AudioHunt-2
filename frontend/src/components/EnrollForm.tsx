import { useEffect, useId, useRef, useState } from 'react';
import { enrollSong } from '../services/api';

interface EnrollFormProps {
  onEnrolled?: () => void;
}

function sanitizeTitle(fileName: string) {
  return fileName.replace(/\.[^.]+$/, '').replace(/[_-]+/g, ' ').trim();
}

export default function EnrollForm({ onEnrolled }: EnrollFormProps) {
  const inputId = useId();
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [title, setTitle] = useState('');
  const [artist, setArtist] = useState('');
  const [audioFile, setAudioFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [status, setStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle');
  const [message, setMessage] = useState('');
  const [uploadHint, setUploadHint] = useState('');

  useEffect(() => {
    if (status !== 'uploading') {
      setUploadHint('');
      return;
    }

    setUploadHint('Uploading and fingerprinting your track...');
    const firstHint = window.setTimeout(() => {
      setUploadHint('Still fingerprinting on the cloud server. Keep this tab open.');
    }, 10000);
    const secondHint = window.setTimeout(() => {
      setUploadHint('Large songs can take around a minute on the free deployment.');
    }, 45000);

    return () => {
      window.clearTimeout(firstHint);
      window.clearTimeout(secondHint);
    };
  }, [status]);

  const assignFile = (file: File | null) => {
    if (!file) {
      return;
    }

    setAudioFile(file);
    setStatus('idle');
    setMessage('');

    if (!title) {
      setTitle(sanitizeTitle(file.name));
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    assignFile(event.target.files?.[0] ?? null);
  };

  const handleDrop = (event: React.DragEvent<HTMLLabelElement>) => {
    event.preventDefault();
    setIsDragging(false);
    assignFile(event.dataTransfer.files?.[0] ?? null);
  };

  const resetForm = () => {
    setTitle('');
    setArtist('');
    setAudioFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!audioFile || !title.trim() || !artist.trim()) {
      setStatus('error');
      setMessage('Add a title, an artist, and one audio file before enrolling.');
      return;
    }

    setStatus('uploading');
    setMessage('');

    try {
      const response = await enrollSong(audioFile, title.trim(), artist.trim());
      setStatus('success');
      setMessage(response.message || `Added to your library as track #${response.id}.`);
      resetForm();
      onEnrolled?.();
    } catch (error: unknown) {
      const responseData = (error as { response?: { data?: unknown } })?.response?.data;
      const nextMessage =
        typeof responseData === 'string'
          ? responseData
          : error instanceof Error
            ? error.message
            : 'Enrollment failed. Please try another file.';
      setStatus('error');
      setMessage(nextMessage);
    }
  };

  return (
    <form className="enroll-form" onSubmit={handleSubmit}>
      <div className="section-heading">
        <div>
          <p className="eyebrow">Library intake</p>
          <h2>Add a song to your catalog</h2>
        </div>
      </div>

      <label
        className={`drop-zone ${isDragging ? 'drop-zone--active' : ''}`}
        htmlFor={inputId}
        onDragEnter={() => setIsDragging(true)}
        onDragLeave={() => setIsDragging(false)}
        onDragOver={(event) => event.preventDefault()}
        onDrop={handleDrop}
      >
        <input
          accept="audio/*"
          id={inputId}
          onChange={handleFileChange}
          ref={fileInputRef}
          type="file"
        />
        <strong>{audioFile ? audioFile.name : 'Drop audio here or browse files'}</strong>
        <span>MP3, WAV, M4A, WEBM or any audio format your backend already accepts.</span>
      </label>

      <div className="field-grid">
        <label className="field">
          <span>Title</span>
          <input
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Shape of You"
            type="text"
            value={title}
          />
        </label>
        <label className="field">
          <span>Artist</span>
          <input
            onChange={(event) => setArtist(event.target.value)}
            placeholder="Ed Sheeran"
            type="text"
            value={artist}
          />
        </label>
      </div>

      <div className="enroll-form__footer">
        <button className="primary-btn" disabled={status === 'uploading'} type="submit">
          {status === 'uploading' ? 'Adding to library...' : 'Enroll track'}
        </button>
        {status === 'uploading' ? <div className="progress-rail" aria-hidden="true" /> : null}
      </div>

      {uploadHint ? (
        <div className="form-message form-message--info">
          {uploadHint}
        </div>
      ) : null}

      {message ? (
        <div className={`form-message form-message--${status === 'success' ? 'success' : 'error'}`}>
          {message}
        </div>
      ) : null}
    </form>
  );
}
