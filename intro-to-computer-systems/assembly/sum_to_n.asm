; Returns the sum of all numbers from 1 to n, inclusive
; TODO bonus points for using math -- ((n ^ 2) + n) / 2 -- to get the answer, rather than a loop

			global 		sum_to_n

			section 	.text
sum_to_n:	mov			rsi, rdi
			cmp			rsi, 0
			jne			.work
			mov			rax, rsi
			ret

.work:		add			rax, rsi
			sub			rsi, 1
			cmp			rsi, 0
			jne			.work
			ret
