; Returns the sum of all numbers from 1 to n, inclusive

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

;; Uses the formula ((n * n) + n) / 2
;; sum_to_n:	mov		rsi, rdi
;; 			imul	rsi, rsi
;; 			add 	rsi, rdi
;; 			mov 	rax, rsi
;; 			mov 	rcx, 2
;; 			div		rcx
;; 			ret
