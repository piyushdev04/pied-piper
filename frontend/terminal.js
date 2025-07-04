const output = document.getElementById('output');
const form = document.getElementById('command-form');
const input = document.getElementById('command-input');
const fileInput = document.getElementById('file-input');

let selectedFile = null;

function print(text) {
  output.innerHTML += text + '\n';
  output.scrollTop = output.scrollHeight;
}

function clearOutput() {
  output.innerHTML = '';
}

function showHelp() {
  print('Commands:');
  print('  select      Choose a file');
  print('  compress    Compress selected file');
  print('  decompress  Decompress .zst file');
  print('  help        Show this help');
  print('  clear       Clear the terminal');
}

form.onsubmit = async (e) => {
  e.preventDefault();
  const cmd = input.value.trim();
  print(`<span class='prompt'>user@piedpiper:~$</span> ${cmd}`);
  input.value = '';

  if (cmd === 'clear') {
    clearOutput();
    return;
  }
  if (cmd === 'help') {
    showHelp();
    return;
  }
  if (cmd === 'select') {
    fileInput.value = '';
    fileInput.click();
    return;
  }
  if (cmd === 'compress') {
    if (!selectedFile) {
      print('No file selected. Use "select" first.');
      return;
    }
    print('Compressing...');
    const formData = new FormData();
    formData.append('file', selectedFile);
    try {
      const res = await fetch('/compress', { method: 'POST', body: formData });
      if (!res.ok) throw new Error('Compression failed');
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = selectedFile.name + '.zst';
      a.click();
      print('Compression complete. Download started.');
    } catch (err) {
      print('Compression error.');
    }
    return;
  }
  if (cmd === 'decompress') {
    if (!selectedFile) {
      print('No file selected. Use "select" first.');
      return;
    }
    print('Decompressing...');
    const formData = new FormData();
    formData.append('file', selectedFile);
    try {
      const res = await fetch('/decompress', { method: 'POST', body: formData });
      if (!res.ok) throw new Error('Decompression failed');
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = selectedFile.name.replace(/\.zst$/, '') || 'decompressed';
      a.click();
      print('Decompression complete. Download started.');
    } catch (err) {
      print('Decompression error.');
    }
    return;
  }
  print('Unknown command. Type "help".');
};

fileInput.onchange = (e) => {
  if (fileInput.files.length > 0) {
    selectedFile = fileInput.files[0];
    print(`Selected file: ${selectedFile.name}`);
  }
};

// Initial help
showHelp(); 