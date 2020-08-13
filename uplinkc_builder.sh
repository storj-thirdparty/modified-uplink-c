#!/bin/bash

grep -n "//export MAKE_CONST=" .build/uplink.h | while read -r line ; do

	line_no=${line%://*}

	updated_line_no=$(($line_no+2))

	list_values=($(echo $line | cut -d "=" -f2 | tr "," "\n"))

	original=$(sed -n "$updated_line_no,$updated_line_no p" .build/uplink.h | cut -d "(" -f2 | cut -d ")" -f1)

	len=$(echo $original | grep -o ',' | wc -l)
	len=$((len+1))

	counter=1
	final=""
	
	updated_array=$original
	IFS=','
	for i in $updated_array
	do
		IFS=''
		if echo ${list_values[@]} | grep -q -w $counter
		then
			# needs const*
			if [[ $counter -eq 1 ]]
			then
				# first element which needs const*
				if [[ $counter -eq $len ]]
				then
					# last element in array
					final="${final}const ${i}"
				else
					# not the last element in array
					final="${final}const ${i},"
				fi
			else
				# non-first element which needs const*
				if [[ $counter -eq $len ]]
				then
					# last element in array
					final="${final} const${i}"
				else
					# not the last element in array
					final="${final} const${i},"
				fi
			fi
		else
			# doesn't need const*
			if [[ $counter -eq 1 ]]
			then
				# first element which doesn't need const*
				if [[ $counter -eq $len ]]
				then
					# last element in array
					final="${final}${i}"
				else
					# not the last element in array
					final="${final}${i},"
				fi
			else
				# non-first element which doesn't need const*
				if [[ $counter -eq $len ]]
				then
					# last element in array
					final="${final}${i}"
				else
					# not the last element in array
					final="${final}${i},"
				fi
			fi

		fi
		
		counter=$((counter+1))
	done

	IFS=''

	original=$(echo "${original}" | sed -e 's/[]$.*[\^]/\\&/g' )
	final=$(echo "${final}" | sed -e 's/[]$.*[\^]/\\&/g' )

	if [[ "$OSTYPE" == "darwin"* ]]
	then
		sed -i '' "${updated_line_no}s/${original}/${final}/g" .build/uplink.h
	else
		sed -i "${updated_line_no}s/${original}/${final}/g" .build/uplink.h
	fi

done

if [[ "$OSTYPE" == "darwin"* ]]
then
	sed -i '' "/MAKE_CONST=/d" .build/uplink.h
else
	sed -i "/MAKE_CONST=/d" .build/uplink.h
fi