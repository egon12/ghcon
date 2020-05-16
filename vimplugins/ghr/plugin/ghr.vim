function! ghr#comment(msg)
	let filePath = expand('%')
	let lineNumber = line('.')
	echom system('ghr comment ' . filePath . ':' . lineNumber . ' "' . a:msg . '"')
endfunction

function! ghr#finish(msg)
	echom system('ghr finish "' . a:msg . '"')
endfunction

function! ghr#cancel()
	echom system('ghr cancel')
endfunction

command -nargs=1 GHRComment call ghr#comment(<args>)

command -nargs=1 GHRFinish call ghr#finish(<args>)

command -nargs=0 GHRCancel call ghr#cancel()
