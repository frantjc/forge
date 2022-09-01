// package contaminate contains ways for forge to internally
// pass state between ores, such as to make sequential ores share
// a filesystem

// this may sound like it breaks the fundamental idea of an ore,
// which is to have an encodable data structure that represents
// one or many containerized commands in their entirety so that
// such commands may be cached and retrieved later by their encoding

// however, this package is only used by individual ores internally
// (which may contain one or many ores), which maintains the ore
// being hermetic

// it is, of course, an internal package so that users cannot do the same,
// thus breaking forge's concept of an ore externally
package contaminate
