package main

import (
	"fmt"
	"strconv"
)

type Command struct {
	Name string
	Args []string
}

type Parser struct {
	data []byte
	pos  int
}

func NewParser(data []byte) *Parser {
	return &Parser{
		data: data,
		pos:  0,
	}
}

func (p *Parser)readLine() ([]byte, error) {
	start := p.pos
	for {
		if p.pos+1 >= len(p.data) {
			return nil, fmt.Errorf("unexpected end of input")
		}
		if p.data[p.pos]=='\r' && p.data[p.pos+1]=='\n'{
			 	line:=p.data[start:p.pos]
				p.pos+=2
				return line,nil
		}
		p.pos++
	}
}

func (p *Parser) readArrayHeader() (int, error) {

	line, err := p.readLine()
	if err != nil {
		return 0, err
	}
	if len(line) == 0 || line[0] != '*' {
		return 0, fmt.Errorf("array expect")
	}

	count, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return 0, fmt.Errorf("invalid array len")
	}

	if count <= 0 {
		return 0, fmt.Errorf("empty command")
	}
	return count, nil
}

func (p *Parser) readBulkString() ([]byte, error) {

	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '$' {
		return nil, fmt.Errorf("expected bulk string")
	}

	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid bulk string length")
	}

	if length < 0 {
		return nil, fmt.Errorf("negative bulk string length")
	}

	if p.pos+length+2 > len(p.data) {
		return nil, fmt.Errorf("incomplete bulk string")
	}

	value := p.data[p.pos : p.pos+length]

	p.pos += length

	if p.data[p.pos] != '\r' || p.data[p.pos+1] != '\n' {
		return nil, fmt.Errorf("missing CRLF after bulk string")
	}

	p.pos += 2

	return value, nil
}

func (p *Parser)Parse() (*Command,error){
	count,err:=p.readArrayHeader()
	if err!=nil{
		return nil,err
	}
	values:=make([]string,0,count)
	
	for i:=0;i<count;i++{
		value,err:=p.readBulkString()
		if err!=nil{
			return nil,err
		}
		values=append(values, string(value))
	}
	cmd:=&Command{
		Name:values[0],
		Args:values[1:],
	}
	return cmd,nil

}