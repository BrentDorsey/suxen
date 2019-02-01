import * as React from 'react'
import Results from './Results'
import gql from "graphql-tag";

import scss from './Search.scss'
import {Query} from 'react-apollo';

const QUERY = gql`
  query search($query: String!) {
    search(query: $query) {
      name
    	version
  		createdAt
    	dockerVersion
    	author
    	operatingSystem
    	preRelease
      checksum {
        sha1
        sha256
      }
      images {
        id
      	path
      	manifestUrl
      	pullUrl
    		version
      	name
      	author
      	createdAt
      }
    }
  }
`

export default class Search extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      query: props.query || '',
    }
  }

  handleInputChange = (e) => {
    this.setState({
      query: e.target.value || '',
    })
    // this.setResultsFromTestData()
  }

  setResultsFromTestData = () => {
    const {query} = this.state
  }

  render() {
    return (
      <div className="container">
        <div className={scss.section}>
          <div className={scss.title}>Sux<div className={scss.suxen}>e</div>n</div>
          <div className={scss.subtitle}>Container viewer for Nexus</div>
          <div className={scss.flex}>
            <input className={scss.input}
                   placeholder='Enter image name to search'
                   value={this.state.query}
                   onChange={this.handleInputChange}/>
            <div className={scss.searchIcon}></div>
          </div>
        </div>
        <Query query={QUERY} variables={{query: this.state.query}}>
          {({loading, error, data}) => {
            if (loading) {
              return (
                <div className="center">
                  <div className="loader"/>
                </div>
              );
            }
            if (error) return `Error! ${error.message}`;

            console.log('Query returned', data, error, loading)

            return (
              <Results imageGroups={data.search}/>
            )
          }}
        </Query>
      </div>
    )
  }
}
